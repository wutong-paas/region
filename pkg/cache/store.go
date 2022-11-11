/*
Copyright 2022 The Wutong-PaaS Authors.

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

package cache

import (
	"sync"

	"github.com/wutong-paas/region/generated/clientset/versioned"
	"github.com/wutong-paas/region/generated/informers/externalversions"
	"github.com/wutong-paas/region/generated/listers/core/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	appsv1 "k8s.io/client-go/listers/apps/v1"
	corev1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	SystemNamespace = "wt-system"
	selector        = "creator=Wutong"
)

type store struct {
	DeploymentLister appsv1.DeploymentLister
	ConfigMapLister  corev1.ConfigMapLister

	SysComponentLister v1alpha1.SysComponentLister
	BizComponentLister v1alpha1.BizComponentLister
}

var cachedStore *store

func Store() *store {
	if cachedStore == nil {
		deposit()

	}
	return cachedStore
}

func deposit() {
	if cachedStore == nil {
		// clients
		clientset := kubernetes.NewForConfigOrDie(ctrl.GetConfigOrDie())
		regionClientset := versioned.NewForConfigOrDie(ctrl.GetConfigOrDie())

		kubeFactory := informers.NewSharedInformerFactoryWithOptions(clientset, 0, informers.WithTweakListOptions(func(options *metav1.ListOptions) {
			options.LabelSelector = selector
		}))

		deployment := kubeFactory.Apps().V1().Deployments()
		deploymentInformer := deployment.Informer()
		configMap := kubeFactory.Core().V1().ConfigMaps()
		configMapInformer := configMap.Informer()
		kubeFactory.Start(wait.NeverStop)

		versionedFactory := externalversions.NewSharedInformerFactory(regionClientset, 0)

		sysComponent := versionedFactory.Core().V1alpha1().SysComponents()
		bizComponent := versionedFactory.Core().V1alpha1().BizComponents()
		sysComponentInformer := sysComponent.Informer()
		bizComponentInformer := bizComponent.Informer()
		versionedFactory.Start(wait.NeverStop)

		sharedInformers := []cache.SharedInformer{
			deploymentInformer,
			configMapInformer,

			sysComponentInformer,
			bizComponentInformer,
		}
		var wg sync.WaitGroup
		wg.Add(len(sharedInformers))
		for _, si := range sharedInformers {
			go func(si cache.SharedInformer) {
				if !cache.WaitForCacheSync(wait.NeverStop, si.HasSynced) {
					panic("timed out waiting for caches to sync")
				}
				wg.Done()
			}(si)
		}
		wg.Wait()

		deploymentLister := deployment.Lister()
		configMapLister := configMap.Lister()

		sysComponentLister := sysComponent.Lister()
		bizComponentList := bizComponent.Lister()

		cachedStore = &store{
			DeploymentLister: deploymentLister,
			ConfigMapLister:  configMapLister,

			SysComponentLister: sysComponentLister,
			BizComponentLister: bizComponentList,
		}
	}
}
