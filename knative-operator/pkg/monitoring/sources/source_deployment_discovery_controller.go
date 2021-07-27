package sources

import (
	"context"
	"os"
	"strconv"

	"github.com/openshift-knative/serverless-operator/knative-operator/pkg/common"
	"github.com/openshift-knative/serverless-operator/knative-operator/pkg/monitoring"
	okomon "github.com/openshift-knative/serverless-operator/openshift-knative-operator/pkg/monitoring"
	v1 "k8s.io/api/apps/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"knative.dev/operator/pkg/apis/operator/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var (
	generateSourceServiceMonitorsEnvVar = "SOURCES_GENERATE_SERVICE_MONITORS"
	useClusterMonitoringEnvVar          = "SOURCES_USE_CLUSTER_MONITORING"
	rbacLabelKey                        = "serverless.monitoring"
	sourceRbacLabels                    = map[string]string{rbacLabelKey: "true"}

	log = common.Log.WithName("source-deployment-discovery-controller")
)

// Add creates a new Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileSourceDeployment{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("source-deployment-discovery-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}
	// common function to enqueue reconcile requests for resources
	enqueueRequests := handler.MapFunc(func(obj client.Object) []reconcile.Request {
		dep := obj.(*v1.Deployment)
		sourceLabel := dep.Spec.Selector.MatchLabels[SourceLabel]
		sourceNameLabel := dep.Spec.Selector.MatchLabels[SourceNameLabel]
		sourceRoleLabel := dep.Spec.Selector.MatchLabels[SourceRoleLabel]

		if (sourceLabel != "" && sourceNameLabel != "") || (sourceLabel != "" && sourceRoleLabel != "") {
			return []reconcile.Request{{
				NamespacedName: types.NamespacedName{Namespace: obj.GetNamespace(), Name: obj.GetName()},
			}}
		}
		return nil
	})

	enqueueRequestsRBAC := handler.MapFunc(func(obj client.Object) []reconcile.Request {
		oLabels := obj.GetLabels()
		if _, ok := oLabels[rbacLabelKey]; ok {
			// Use a special request to trigger ns level reconciliation related to sources
			return []reconcile.Request{{
				NamespacedName: types.NamespacedName{Namespace: obj.GetNamespace(), Name: ""},
			}}
		}
		return nil
	})

	if err = c.Watch(&source.Kind{Type: &v1.Deployment{}}, handler.EnqueueRequestsFromMapFunc(enqueueRequests), skipDeletePredicate{}); err != nil {
		return err
	}
	if err = c.Watch(&source.Kind{Type: &rbacv1.RoleBinding{}}, handler.EnqueueRequestsFromMapFunc(enqueueRequestsRBAC), skipCreatePredicate{}); err != nil {
		return err
	}
	return c.Watch(&source.Kind{Type: &rbacv1.Role{}}, handler.EnqueueRequestsFromMapFunc(enqueueRequestsRBAC), skipCreatePredicate{})
}

// blank assignment to verify that ReconcileSourceDeployment implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileSourceDeployment{}

// ReconcileSourceDeployment reconciles a source deployment object
type ReconcileSourceDeployment struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for an eventing source deployment
func (r *ReconcileSourceDeployment) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	reconcileForSource := true
	if request.Name == "" {
		// We will handle only rbac resources this time
		reconcileForSource = false
	}
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling the source deployment, setting up a service/service monitor if required")
	dep := &v1.Deployment{}
	inDeletion := false
	if reconcileForSource {
		err := r.client.Get(context.TODO(), request.NamespacedName, dep)
		if apierrors.IsNotFound(err) {
			// The deployment does not exist anymore, deletions are shown as failing reads
			log.Info("Source in deletion phase")
			inDeletion = true
		} else if err != nil {
			return reconcile.Result{}, err
		}
	}
	eventing := &v1alpha1.KnativeEventing{}
	if err := r.client.Get(context.TODO(), types.NamespacedName{Namespace: "knative-eventing", Name: "knative-eventing"}, eventing); err != nil {
		return reconcile.Result{}, err
	}
	// If monitoring is set to on/off this triggers a global resync to source adapters.
	// Same applies if we change any of the env vars affecting cluster monitoring or service monitor resource generation.
	// The Serverless operator pod is restarted and local informer caches are synched.
	if okomon.ShouldEnableMonitoring(eventing.Spec.GetConfig()) {
		// If in deletion there is nothing to be done, owner refs will remove source service monitors
		// Make sure we do not setup any resources if the source is being deleted
		// A deletion event will make sure that we detect a deletion properly from cluster state
		if !inDeletion {
			if err := r.setupClusterMonitoringForSources(request.Namespace); err != nil {
				return reconcile.Result{}, err
			}
			if reconcileForSource {
				if err := r.generateSourceServiceMonitors(dep); err != nil {
					return reconcile.Result{}, err
				}
			}
		}
	} else {
		// Remove any relics if previously monitoring was on.
		if dep.Namespace != "knative-eventing" {
			if err := monitoring.RemoveClusterMonitoringRequirements(r.client, nil, dep.GetNamespace(), sourceRbacLabels); err != nil {
				return reconcile.Result{}, err
			}
			if reconcileForSource {
				if err := RemoveSourceServiceMonitorResources(r.client, dep); err != nil {
					return reconcile.Result{}, err
				}
			}
		}
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileSourceDeployment) generateSourceServiceMonitors(dep *v1.Deployment) error {
	shouldGenerateSourceMonitors, err := r.shouldGenerateSourceServiceMonitorsByDefault()
	if err != nil {
		return err
	}
	if shouldGenerateSourceMonitors {
		if err := SetupSourceServiceMonitorResources(r.client, dep); err != nil {
			return err
		}
	} else {
		if dep.Namespace != "knative-eventing" {
			if err := RemoveSourceServiceMonitorResources(r.client, dep); err != nil {
				return err
			}
		}
	}
	return nil
}

// Setup cluster monitoring Prometheus monitoring requirements
func (r *ReconcileSourceDeployment) setupClusterMonitoringForSources(ns string) error {
	shouldEnableClusterMonitoring, err := r.shouldUseClusterMonitoringForSourcesByDefault()
	if err != nil {
		return err
	}
	if shouldEnableClusterMonitoring {
		if err := monitoring.SetupClusterMonitoringRequirements(r.client, nil, ns, sourceRbacLabels); err != nil {
			return err
		}
	} else {
		// Make sure we disable cluster monitoring if we have to eg. we move from a state of enabled to disabled and
		// resources are left without cleanup. This brings us to the right state.
		if ns != "knative-eventing" {
			if err := monitoring.RemoveClusterMonitoringRequirements(r.client, nil, ns, sourceRbacLabels); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *ReconcileSourceDeployment) shouldUseClusterMonitoringForSourcesByDefault() (bool, error) {
	enable, err := strconv.ParseBool(os.Getenv(useClusterMonitoringEnvVar))
	return enable, err
}

func (r *ReconcileSourceDeployment) shouldGenerateSourceServiceMonitorsByDefault() (bool, error) {
	enable, err := strconv.ParseBool(os.Getenv(generateSourceServiceMonitorsEnvVar))
	return enable, err
}

type skipCreatePredicate struct {
	predicate.Funcs
}

func (skipCreatePredicate) Create(e event.CreateEvent) bool {
	return false
}

type skipDeletePredicate struct {
	predicate.Funcs
}

func (skipDeletePredicate) Delete(e event.DeleteEvent) bool {
	return false
}
