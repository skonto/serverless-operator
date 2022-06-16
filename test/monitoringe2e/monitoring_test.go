package monitoringe2e

import (
	"github.com/openshift-knative/serverless-operator/knative-operator/pkg/monitoring"
	"testing"

	"github.com/openshift-knative/serverless-operator/test"
)

func TestKnativeMetrics(t *testing.T) {
	caCtx := test.SetupClusterAdmin(t)
	cleanup := func() {
		test.CleanupAll(t, caCtx)
	}
	test.CleanupOnInterrupt(t, cleanup)
	defer cleanup()
	t.Setenv(monitoring.OperatorDeploymentNameEnvKey, "knative-openshift")
	if err := SetUpServerlessOperatorMonitoring(caCtx); err != nil {
		t.Fatal("Failed to setup Serverless Operator requirements", err)
	}

	t.Run("verify Serving control plane metrics work correctly", func(t *testing.T) {
		// Serving control plane metrics should work
		if err := VerifyMetrics(caCtx, servingMetricQueries); err != nil {
			t.Fatal("Failed to verify that Serving control plane metrics work correctly", err)
		}
	})

	t.Run("verify Eventing control plane metrics work correctly", func(t *testing.T) {
		// Eventing control plane metrics should work
		if err := VerifyMetrics(caCtx, eventingMetricQueries); err != nil {
			t.Fatal("Failed to verify that Eventing control plane metrics work correctly", err)
		}
	})

	t.Run("verify Knative operators and Openshift ingress metrics work correctly", func(t *testing.T) {
		// Knative operators and Openshift ingress metrics should work
		if err := VerifyMetrics(caCtx, serverlessComponentQueries); err != nil {
			t.Fatal("Failed to verify that Knative operators and Openshift ingress metrics work correctly", err)
		}
	})
}
