/*
Copyright The Kubernetes Authors.

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

package v1alpha1

import (
	"time"

	v1alpha1 "github.com/tacf/k8s-pingdom-operator/pkg/apis/pingdomcontroller/v1alpha1"
	scheme "github.com/tacf/k8s-pingdom-operator/pkg/generated/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// PingdomOperatorsGetter has a method to return a PingdomOperatorInterface.
// A group's client should implement this interface.
type PingdomOperatorsGetter interface {
	PingdomOperators(namespace string) PingdomOperatorInterface
}

// PingdomOperatorInterface has methods to work with PingdomOperator resources.
type PingdomOperatorInterface interface {
	Create(*v1alpha1.PingdomOperator) (*v1alpha1.PingdomOperator, error)
	Update(*v1alpha1.PingdomOperator) (*v1alpha1.PingdomOperator, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.PingdomOperator, error)
	List(opts v1.ListOptions) (*v1alpha1.PingdomOperatorList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.PingdomOperator, err error)
	PingdomOperatorExpansion
}

// pingdomOperators implements PingdomOperatorInterface
type pingdomOperators struct {
	client rest.Interface
	ns     string
}

// newPingdomOperators returns a PingdomOperators
func newPingdomOperators(c *PingdomcontrollerV1alpha1Client, namespace string) *pingdomOperators {
	return &pingdomOperators{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the pingdomOperator, and returns the corresponding pingdomOperator object, and an error if there is any.
func (c *pingdomOperators) Get(name string, options v1.GetOptions) (result *v1alpha1.PingdomOperator, err error) {
	result = &v1alpha1.PingdomOperator{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("pingdomoperators").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of PingdomOperators that match those selectors.
func (c *pingdomOperators) List(opts v1.ListOptions) (result *v1alpha1.PingdomOperatorList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.PingdomOperatorList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("pingdomoperators").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested pingdomOperators.
func (c *pingdomOperators) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("pingdomoperators").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a pingdomOperator and creates it.  Returns the server's representation of the pingdomOperator, and an error, if there is any.
func (c *pingdomOperators) Create(pingdomOperator *v1alpha1.PingdomOperator) (result *v1alpha1.PingdomOperator, err error) {
	result = &v1alpha1.PingdomOperator{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("pingdomoperators").
		Body(pingdomOperator).
		Do().
		Into(result)
	return
}

// Update takes the representation of a pingdomOperator and updates it. Returns the server's representation of the pingdomOperator, and an error, if there is any.
func (c *pingdomOperators) Update(pingdomOperator *v1alpha1.PingdomOperator) (result *v1alpha1.PingdomOperator, err error) {
	result = &v1alpha1.PingdomOperator{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("pingdomoperators").
		Name(pingdomOperator.Name).
		Body(pingdomOperator).
		Do().
		Into(result)
	return
}

// Delete takes name of the pingdomOperator and deletes it. Returns an error if one occurs.
func (c *pingdomOperators) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("pingdomoperators").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *pingdomOperators) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("pingdomoperators").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched pingdomOperator.
func (c *pingdomOperators) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.PingdomOperator, err error) {
	result = &v1alpha1.PingdomOperator{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("pingdomoperators").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}