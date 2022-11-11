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

// SysComponentLister helps list SysComponents.
// All objects returned here must be treated as read-only.
type SysComponentLister interface {
	// List lists all SysComponents in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.SysComponent, err error)
	// SysComponents returns an object that can list and get SysComponents.
	SysComponents(namespace string) SysComponentNamespaceLister
	SysComponentListerExpansion
}

// sysComponentLister implements the SysComponentLister interface.
type sysComponentLister struct {
	indexer cache.Indexer
}

// NewSysComponentLister returns a new SysComponentLister.
func NewSysComponentLister(indexer cache.Indexer) SysComponentLister {
	return &sysComponentLister{indexer: indexer}
}

// List lists all SysComponents in the indexer.
func (s *sysComponentLister) List(selector labels.Selector) (ret []*v1alpha1.SysComponent, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.SysComponent))
	})
	return ret, err
}

// SysComponents returns an object that can list and get SysComponents.
func (s *sysComponentLister) SysComponents(namespace string) SysComponentNamespaceLister {
	return sysComponentNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// SysComponentNamespaceLister helps list and get SysComponents.
// All objects returned here must be treated as read-only.
type SysComponentNamespaceLister interface {
	// List lists all SysComponents in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.SysComponent, err error)
	// Get retrieves the SysComponent from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.SysComponent, error)
	SysComponentNamespaceListerExpansion
}

// sysComponentNamespaceLister implements the SysComponentNamespaceLister
// interface.
type sysComponentNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all SysComponents in the indexer for a given namespace.
func (s sysComponentNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.SysComponent, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.SysComponent))
	})
	return ret, err
}

// Get retrieves the SysComponent from the indexer for a given namespace and name.
func (s sysComponentNamespaceLister) Get(name string) (*v1alpha1.SysComponent, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("syscomponent"), name)
	}
	return obj.(*v1alpha1.SysComponent), nil
}