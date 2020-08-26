package common

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	knativeUp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "knative_up",
			Help:        "Reports if a Knative component is up",
			ConstLabels: map[string]string{"type": "serving_status"},
		},
		[]string{"type"},
		)
	KnativeServingUpG = knativeUp.WithLabelValues("serving_status")
	KnativeEventingUpG = knativeUp.WithLabelValues("eventing_status")
)

func init() {
	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(KnativeServingUpG, KnativeEventingUpG)
}
