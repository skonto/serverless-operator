package common

import (
	"context"
	"fmt"

	"k8s.io/client-go/kubernetes/scheme"
	mfclient "github.com/manifestival/controller-runtime-client"
	mf "github.com/manifestival/manifestival"
	"k8s.io/apimachinery/pkg/runtime/schema"
	eventingv1alpha1 "knative.dev/operator/pkg/apis/operator/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	kubemetrics "github.com/operator-framework/operator-sdk/pkg/kube-metrics"
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/operator-framework/operator-sdk/pkg/metrics"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)
const EventingBrokerServiceMoinitorPath = "deploy/resources/broker-service-monitors.yaml"

func SetupServerlessOperatorServiceMonitor(ctx context.Context, cfg *rest.Config, api client.Client, metricsPort int32, metricsHost string, operatorMetricsPort int32) error {
	// Commented below to avoid a stream of these errors at startup:
	// E1021 22:50:03.372487       1 reflector.go:134] github.com/operator-framework/operator-sdk/pkg/kube-metrics/collector.go:67: Failed to list *unstructured.Unstructured: the server could not find the requested resource
	if err := serveCRMetrics(cfg, metricsHost, operatorMetricsPort); err != nil {
		log.Info("Could not generate and serve custom resource metrics", "error", err.Error())
	}

	// Add to the below struct any other metrics ports you want to expose.
	servicePorts := []v1.ServicePort{
		{Port: metricsPort, Name: metrics.OperatorPortName, Protocol: v1.ProtocolTCP, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: metricsPort}},
		{Port: operatorMetricsPort, Name: metrics.CRPortName, Protocol: v1.ProtocolTCP, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: operatorMetricsPort}},
	}
	// Create Service object to expose the metrics port(s).
	service, err := metrics.CreateMetricsService(ctx, cfg, servicePorts)
	if err != nil {
		log.Info("Could not create metrics Service", "error", err.Error())
	}

	// CreateServiceMonitors will automatically create the prometheus-operator ServiceMonitor resources
	// necessary to configure Prometheus to scrape metrics from this operator.
	services := []*v1.Service{service}
	metricsNamespace, err := k8sutil.GetOperatorNamespace()
	if err != nil {
		log.Error(err, "failed to get metrics namespace")
		return err
	}
	_, err = metrics.CreateServiceMonitors(cfg, metricsNamespace, services)

	if err != nil {
		log.Info("Could not create ServiceMonitor object", "error", err.Error())
		// If this operator is deployed to a cluster without the prometheus-operator running, it will return
		// ErrServiceMonitorNotPresent, which can be used to safely skip ServiceMonitor creation.
		if err == metrics.ErrServiceMonitorNotPresent {
			log.Info("Install prometheus-operator in your cluster to create ServiceMonitor objects", "error", err.Error())
		}
	}
	return err
}

// serveCRMetrics gets the Operator/CustomResource GVKs and generates metrics based on those types.
// It serves those metrics on "http://metricsHost:operatorMetricsPort".
func serveCRMetrics(cfg *rest.Config, metricsHost string, operatorMetricsPort int32) error {
	gvkFilterList := []schema.GroupVersionKind{
		schema.GroupVersionKind{
			Group:   "operator.knative.dev",
			Version: "v1alpha1",
			Kind:    "KnativeServing",
		},
		schema.GroupVersionKind{
			Group: "operator.knative.dev",
			Version: "v1alpha1",
			Kind: "KnativeEventing",
		},
	}
	// To generate metrics in other namespaces, add the values below.
	ns := []string{""}
	// Generate and serve custom resource specific metrics.
	err := kubemetrics.GenerateAndServeCRMetrics(cfg, ns, gvkFilterList, metricsHost, operatorMetricsPort)
	if err != nil {
		return err
	}
	return nil
}

func SetupEventingServiceMonitors(namespace string, instance *eventingv1alpha1.KnativeEventing)  error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return fmt.Errorf("failed to create cluster config: %w", err)
	}
	cl, err := client.New(config, client.Options{})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	manifest, err := mf.NewManifest(EventingBrokerServiceMoinitorPath, mf.UseClient(mfclient.NewClient(cl)))
	if err != nil {
		return fmt.Errorf("unable to parse broker services: %w", err)
	}
	transforms := []mf.Transformer{mf.InjectOwner(instance), mf.InjectNamespace(namespace)}
	if manifest, err = manifest.Transform(transforms...); err != nil {
		return fmt.Errorf("unable to transform broker service monitors manifest: %w", err)
	}
	mon := &monitoringv1.ServiceMonitor{}

	var SchemeGroupVersion = schema.GroupVersion{Group: "monitoring.coreos.com", Version: "v1"}
	scheme.Scheme.AddKnownTypes(SchemeGroupVersion, mon)

	//if err := scheme.Scheme.Convert(&manifest.Resources()[0], mon, nil); err != nil {
	//	return err
	//}
	if err := manifest.Apply(); err != nil {
		//errr := cl.Create(context.Background(), mon)
		//if errr != nil {
		//	return fmt.Errorf("unable to create broker service monitors second %w ---> %w", errr, err)
		//}
		//return fmt.Errorf("unable to create broker service monitors %w $$$$$ %v", err, manifest)
		return err
	}
	return nil
}

