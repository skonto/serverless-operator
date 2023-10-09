/*
Copyright 2021 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1beta2

import (
	"context"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	v1beta2 "knative.dev/eventing/pkg/apis/eventing/v1beta2"
	scheme "knative.dev/eventing/pkg/client/clientset/versioned/scheme"
)

// EventTypesGetter has a method to return a EventTypeInterface.
// A group's client should implement this interface.
type EventTypesGetter interface {
	EventTypes(namespace string) EventTypeInterface
}

// EventTypeInterface has methods to work with EventType resources.
type EventTypeInterface interface {
	Create(ctx context.Context, eventType *v1beta2.EventType, opts v1.CreateOptions) (*v1beta2.EventType, error)
	Update(ctx context.Context, eventType *v1beta2.EventType, opts v1.UpdateOptions) (*v1beta2.EventType, error)
	UpdateStatus(ctx context.Context, eventType *v1beta2.EventType, opts v1.UpdateOptions) (*v1beta2.EventType, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta2.EventType, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1beta2.EventTypeList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta2.EventType, err error)
	EventTypeExpansion
}

// eventTypes implements EventTypeInterface
type eventTypes struct {
	client rest.Interface
	ns     string
}

// newEventTypes returns a EventTypes
func newEventTypes(c *EventingV1beta2Client, namespace string) *eventTypes {
	return &eventTypes{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the eventType, and returns the corresponding eventType object, and an error if there is any.
func (c *eventTypes) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta2.EventType, err error) {
	result = &v1beta2.EventType{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("eventtypes").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of EventTypes that match those selectors.
func (c *eventTypes) List(ctx context.Context, opts v1.ListOptions) (result *v1beta2.EventTypeList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1beta2.EventTypeList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("eventtypes").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested eventTypes.
func (c *eventTypes) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("eventtypes").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a eventType and creates it.  Returns the server's representation of the eventType, and an error, if there is any.
func (c *eventTypes) Create(ctx context.Context, eventType *v1beta2.EventType, opts v1.CreateOptions) (result *v1beta2.EventType, err error) {
	result = &v1beta2.EventType{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("eventtypes").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(eventType).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a eventType and updates it. Returns the server's representation of the eventType, and an error, if there is any.
func (c *eventTypes) Update(ctx context.Context, eventType *v1beta2.EventType, opts v1.UpdateOptions) (result *v1beta2.EventType, err error) {
	result = &v1beta2.EventType{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("eventtypes").
		Name(eventType.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(eventType).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *eventTypes) UpdateStatus(ctx context.Context, eventType *v1beta2.EventType, opts v1.UpdateOptions) (result *v1beta2.EventType, err error) {
	result = &v1beta2.EventType{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("eventtypes").
		Name(eventType.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(eventType).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the eventType and deletes it. Returns an error if one occurs.
func (c *eventTypes) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("eventtypes").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *eventTypes) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("eventtypes").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched eventType.
func (c *eventTypes) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta2.EventType, err error) {
	result = &v1beta2.EventType{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("eventtypes").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
