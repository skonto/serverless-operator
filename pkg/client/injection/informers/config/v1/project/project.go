// Code generated by injection-gen. DO NOT EDIT.

package project

import (
	context "context"

	v1 "github.com/openshift-knative/serverless-operator/pkg/client/informers/externalversions/config/v1"
	factory "github.com/openshift-knative/serverless-operator/pkg/client/injection/informers/factory"
	controller "knative.dev/pkg/controller"
	injection "knative.dev/pkg/injection"
	logging "knative.dev/pkg/logging"
)

func init() {
	injection.Default.RegisterInformer(withInformer)
}

// Key is used for associating the Informer inside the context.Context.
type Key struct{}

func withInformer(ctx context.Context) (context.Context, controller.Informer) {
	f := factory.Get(ctx)
	inf := f.Config().V1().Projects()
	return context.WithValue(ctx, Key{}, inf), inf.Informer()
}

// Get extracts the typed informer from the context.
func Get(ctx context.Context) v1.ProjectInformer {
	untyped := ctx.Value(Key{})
	if untyped == nil {
		logging.FromContext(ctx).Panic(
			"Unable to fetch github.com/openshift-knative/serverless-operator/pkg/client/informers/externalversions/config/v1.ProjectInformer from context.")
	}
	return untyped.(v1.ProjectInformer)
}
