// Code generated by injection-gen. DO NOT EDIT.

package clusterversion

import (
	context "context"

	versioned "github.com/openshift-knative/serverless-operator/pkg/client/clientset/versioned"
	v1 "github.com/openshift-knative/serverless-operator/pkg/client/informers/externalversions/config/v1"
	client "github.com/openshift-knative/serverless-operator/pkg/client/injection/client"
	factory "github.com/openshift-knative/serverless-operator/pkg/client/injection/informers/factory"
	configv1 "github.com/openshift-knative/serverless-operator/pkg/client/listers/config/v1"
	apiconfigv1 "github.com/openshift/api/config/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	cache "k8s.io/client-go/tools/cache"
	controller "knative.dev/pkg/controller"
	injection "knative.dev/pkg/injection"
	logging "knative.dev/pkg/logging"
)

func init() {
	injection.Default.RegisterInformer(withInformer)
	injection.Dynamic.RegisterDynamicInformer(withDynamicInformer)
}

// Key is used for associating the Informer inside the context.Context.
type Key struct{}

func withInformer(ctx context.Context) (context.Context, controller.Informer) {
	f := factory.Get(ctx)
	inf := f.Config().V1().ClusterVersions()
	return context.WithValue(ctx, Key{}, inf), inf.Informer()
}

func withDynamicInformer(ctx context.Context) context.Context {
	inf := &wrapper{client: client.Get(ctx), resourceVersion: injection.GetResourceVersion(ctx)}
	return context.WithValue(ctx, Key{}, inf)
}

// Get extracts the typed informer from the context.
func Get(ctx context.Context) v1.ClusterVersionInformer {
	untyped := ctx.Value(Key{})
	if untyped == nil {
		logging.FromContext(ctx).Panic(
			"Unable to fetch github.com/openshift-knative/serverless-operator/pkg/client/informers/externalversions/config/v1.ClusterVersionInformer from context.")
	}
	return untyped.(v1.ClusterVersionInformer)
}

type wrapper struct {
	client versioned.Interface

	resourceVersion string
}

var _ v1.ClusterVersionInformer = (*wrapper)(nil)
var _ configv1.ClusterVersionLister = (*wrapper)(nil)

func (w *wrapper) Informer() cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(nil, &apiconfigv1.ClusterVersion{}, 0, nil)
}

func (w *wrapper) Lister() configv1.ClusterVersionLister {
	return w
}

// SetResourceVersion allows consumers to adjust the minimum resourceVersion
// used by the underlying client.  It is not accessible via the standard
// lister interface, but can be accessed through a user-defined interface and
// an implementation check e.g. rvs, ok := foo.(ResourceVersionSetter)
func (w *wrapper) SetResourceVersion(resourceVersion string) {
	w.resourceVersion = resourceVersion
}

func (w *wrapper) List(selector labels.Selector) (ret []*apiconfigv1.ClusterVersion, err error) {
	lo, err := w.client.ConfigV1().ClusterVersions().List(context.TODO(), metav1.ListOptions{
		LabelSelector:   selector.String(),
		ResourceVersion: w.resourceVersion,
	})
	if err != nil {
		return nil, err
	}
	for idx := range lo.Items {
		ret = append(ret, &lo.Items[idx])
	}
	return ret, nil
}

func (w *wrapper) Get(name string) (*apiconfigv1.ClusterVersion, error) {
	return w.client.ConfigV1().ClusterVersions().Get(context.TODO(), name, metav1.GetOptions{
		ResourceVersion: w.resourceVersion,
	})
}
