/*
Copyright 2022 The Knative Authors

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

// Code generated by lister-gen. DO NOT EDIT.

package v1beta1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	v1beta1 "knative.dev/operator/pkg/apis/operator/v1beta1"
)

// KnativeEventingLister helps list KnativeEventings.
// All objects returned here must be treated as read-only.
type KnativeEventingLister interface {
	// List lists all KnativeEventings in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta1.KnativeEventing, err error)
	// KnativeEventings returns an object that can list and get KnativeEventings.
	KnativeEventings(namespace string) KnativeEventingNamespaceLister
	KnativeEventingListerExpansion
}

// knativeEventingLister implements the KnativeEventingLister interface.
type knativeEventingLister struct {
	indexer cache.Indexer
}

// NewKnativeEventingLister returns a new KnativeEventingLister.
func NewKnativeEventingLister(indexer cache.Indexer) KnativeEventingLister {
	return &knativeEventingLister{indexer: indexer}
}

// List lists all KnativeEventings in the indexer.
func (s *knativeEventingLister) List(selector labels.Selector) (ret []*v1beta1.KnativeEventing, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.KnativeEventing))
	})
	return ret, err
}

// KnativeEventings returns an object that can list and get KnativeEventings.
func (s *knativeEventingLister) KnativeEventings(namespace string) KnativeEventingNamespaceLister {
	return knativeEventingNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// KnativeEventingNamespaceLister helps list and get KnativeEventings.
// All objects returned here must be treated as read-only.
type KnativeEventingNamespaceLister interface {
	// List lists all KnativeEventings in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta1.KnativeEventing, err error)
	// Get retrieves the KnativeEventing from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1beta1.KnativeEventing, error)
	KnativeEventingNamespaceListerExpansion
}

// knativeEventingNamespaceLister implements the KnativeEventingNamespaceLister
// interface.
type knativeEventingNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all KnativeEventings in the indexer for a given namespace.
func (s knativeEventingNamespaceLister) List(selector labels.Selector) (ret []*v1beta1.KnativeEventing, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.KnativeEventing))
	})
	return ret, err
}

// Get retrieves the KnativeEventing from the indexer for a given namespace and name.
func (s knativeEventingNamespaceLister) Get(name string) (*v1beta1.KnativeEventing, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta1.Resource("knativeeventing"), name)
	}
	return obj.(*v1beta1.KnativeEventing), nil
}
