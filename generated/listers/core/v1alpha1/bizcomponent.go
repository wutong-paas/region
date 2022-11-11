/*
Copyright 2022 The Wutong Authors.

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

package v1alpha1

import (
	v1alpha1 "github.com/wutong-paas/region/apis/core/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// BizComponentLister helps list BizComponents.
// All objects returned here must be treated as read-only.
type BizComponentLister interface {
	// List lists all BizComponents in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.BizComponent, err error)
	// BizComponents returns an object that can list and get BizComponents.
	BizComponents(namespace string) BizComponentNamespaceLister
	BizComponentListerExpansion
}

// bizComponentLister implements the BizComponentLister interface.
type bizComponentLister struct {
	indexer cache.Indexer
}

// NewBizComponentLister returns a new BizComponentLister.
func NewBizComponentLister(indexer cache.Indexer) BizComponentLister {
	return &bizComponentLister{indexer: indexer}
}

// List lists all BizComponents in the indexer.
func (s *bizComponentLister) List(selector labels.Selector) (ret []*v1alpha1.BizComponent, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.BizComponent))
	})
	return ret, err
}

// BizComponents returns an object that can list and get BizComponents.
func (s *bizComponentLister) BizComponents(namespace string) BizComponentNamespaceLister {
	return bizComponentNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// BizComponentNamespaceLister helps list and get BizComponents.
// All objects returned here must be treated as read-only.
type BizComponentNamespaceLister interface {
	// List lists all BizComponents in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.BizComponent, err error)
	// Get retrieves the BizComponent from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.BizComponent, error)
	BizComponentNamespaceListerExpansion
}

// bizComponentNamespaceLister implements the BizComponentNamespaceLister
// interface.
type bizComponentNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all BizComponents in the indexer for a given namespace.
func (s bizComponentNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.BizComponent, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.BizComponent))
	})
	return ret, err
}

// Get retrieves the BizComponent from the indexer for a given namespace and name.
func (s bizComponentNamespaceLister) Get(name string) (*v1alpha1.BizComponent, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("bizcomponent"), name)
	}
	return obj.(*v1alpha1.BizComponent), nil
}
