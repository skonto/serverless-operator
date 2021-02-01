// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	scheme "github.com/openshift-knative/serverless-operator/pkg/client/clientset/versioned/scheme"
	v1 "github.com/openshift/api/config/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ProxiesGetter has a method to return a ProxyInterface.
// A group's client should implement this interface.
type ProxiesGetter interface {
	Proxies() ProxyInterface
}

// ProxyInterface has methods to work with Proxy resources.
type ProxyInterface interface {
	Create(ctx context.Context, proxy *v1.Proxy, opts metav1.CreateOptions) (*v1.Proxy, error)
	Update(ctx context.Context, proxy *v1.Proxy, opts metav1.UpdateOptions) (*v1.Proxy, error)
	UpdateStatus(ctx context.Context, proxy *v1.Proxy, opts metav1.UpdateOptions) (*v1.Proxy, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Proxy, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.ProxyList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Proxy, err error)
	ProxyExpansion
}

// proxies implements ProxyInterface
type proxies struct {
	client rest.Interface
}

// newProxies returns a Proxies
func newProxies(c *ConfigV1Client) *proxies {
	return &proxies{
		client: c.RESTClient(),
	}
}

// Get takes name of the proxy, and returns the corresponding proxy object, and an error if there is any.
func (c *proxies) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.Proxy, err error) {
	result = &v1.Proxy{}
	err = c.client.Get().
		Resource("proxies").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Proxies that match those selectors.
func (c *proxies) List(ctx context.Context, opts metav1.ListOptions) (result *v1.ProxyList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.ProxyList{}
	err = c.client.Get().
		Resource("proxies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested proxies.
func (c *proxies) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("proxies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a proxy and creates it.  Returns the server's representation of the proxy, and an error, if there is any.
func (c *proxies) Create(ctx context.Context, proxy *v1.Proxy, opts metav1.CreateOptions) (result *v1.Proxy, err error) {
	result = &v1.Proxy{}
	err = c.client.Post().
		Resource("proxies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(proxy).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a proxy and updates it. Returns the server's representation of the proxy, and an error, if there is any.
func (c *proxies) Update(ctx context.Context, proxy *v1.Proxy, opts metav1.UpdateOptions) (result *v1.Proxy, err error) {
	result = &v1.Proxy{}
	err = c.client.Put().
		Resource("proxies").
		Name(proxy.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(proxy).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *proxies) UpdateStatus(ctx context.Context, proxy *v1.Proxy, opts metav1.UpdateOptions) (result *v1.Proxy, err error) {
	result = &v1.Proxy{}
	err = c.client.Put().
		Resource("proxies").
		Name(proxy.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(proxy).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the proxy and deletes it. Returns an error if one occurs.
func (c *proxies) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("proxies").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *proxies) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("proxies").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched proxy.
func (c *proxies) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Proxy, err error) {
	result = &v1.Proxy{}
	err = c.client.Patch(pt).
		Resource("proxies").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}