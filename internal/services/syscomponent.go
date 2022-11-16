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

package services

import (
	"errors"

	"github.com/wutong-paas/region/apis/core/v1alpha1"
	"github.com/wutong-paas/region/pkg/cache"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog/v2"
)

func (s CommonServicer) ListSysComponents() (map[string]*SysComponent, error) {
	sysCompConfigs, err := CachedSysComponentConfigs()
	if err != nil {
		return nil, err
	}

	res := make(map[string]*SysComponent)

	sysComponents, err := cache.Store().SysComponentLister.SysComponents(metav1.NamespaceAll).List(labels.Everything())
	if err != nil {
		klog.Errorf("list sys components failed: %v", err)
		return nil, err
	}

	for k, conf := range sysCompConfigs {
		availableVersions := make([]string, 0, len(conf.AvailableVersions))
		for k := range conf.AvailableVersions {
			availableVersions = append(availableVersions, k)
		}

		res[k] = &SysComponent{
			Name:              k,
			Namespace:         conf.Namespace,
			Description:       conf.Description,
			Status:            SysComponentUnInstalled,
			AvailableVersions: availableVersions,
		}
	}

	for _, sc := range sysComponents {
		item, ok := res[sc.Name]
		if ok {
			item.Status = SysComponentStatus(sc.Status.Phase)
			item.CurrentVersion = sc.Spec.CurrentVersion
		}
	}

	return res, nil
}

func (s CommonServicer) InstallComponnet(name, version string) error {
	sysCompConfigs, err := CachedSysComponentConfigs()
	if err != nil {
		return err
	}

	target, ok := sysCompConfigs[name]
	if !ok {
		return errors.New("sys component not found")
	}

	_, err = s.RegionClient.CoreV1alpha1().SysComponents(target.Namespace).Create(context.TODO(), &v1alpha1.SysComponent{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: target.Namespace,
			Labels: map[string]string{
				"creator": "Wutong",
			},
		},
		Spec: v1alpha1.SysComponentSpec{
			Description:       target.Description,
			AvailableVersions: target.AvailableVersions,
			CurrentVersion:    version,
			InstallWay:        target.InstallWay,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (s CommonServicer) UninstallComponnet(name string) error {
	sysCompConfigs, err := CachedSysComponentConfigs()
	if err != nil {
		return err
	}

	target, ok := sysCompConfigs[name]
	if !ok {
		return errors.New("sys component not found")
	}
	comp, err := s.RegionClient.CoreV1alpha1().SysComponents(target.Namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	comp.Status.Phase = v1alpha1.SysComponentUnInstalling
	_, err = s.RegionClient.CoreV1alpha1().SysComponents(target.Namespace).UpdateStatus(context.TODO(), comp, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (s CommonServicer) UpgradeComponnet(name, version string) error {
	sysCompConfigs, err := CachedSysComponentConfigs()
	if err != nil {
		return err
	}

	target, ok := sysCompConfigs[name]
	if !ok {
		return errors.New("sys component not found")
	}
	comp, err := s.RegionClient.CoreV1alpha1().SysComponents(target.Namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	comp.Spec.CurrentVersion = version
	_, err = s.RegionClient.CoreV1alpha1().SysComponents(target.Namespace).Update(context.TODO(), comp, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}
