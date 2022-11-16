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
	"github.com/wutong-paas/region/pkg/cache"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/yaml"
)

const (
	wtSysComponentConfigMapName = "wt-syscomponent-config"
)

var (
	cachedSysComponentConfigMap *v1.ConfigMap
	cachedSysComponentConfigs   map[string]*SysComponentConfig
)

func CachedSysComponentConfigs() (map[string]*SysComponentConfig, error) {
	// var cachedSysComponentConfigs map[string]*SysComponentConfig
	sysCompConfigMap, err := cache.Store().ConfigMapLister.ConfigMaps(cache.SystemNamespace).Get(wtSysComponentConfigMapName)
	if err != nil {
		klog.Errorf("failed to get syscomponents configmap: %v", err)
		return cachedSysComponentConfigs, err
	}
	if len(cachedSysComponentConfigs) == 0 || cachedSysComponentConfigMap.ResourceVersion != sysCompConfigMap.ResourceVersion {
		cachedSysComponentConfigMap = sysCompConfigMap
		content := sysCompConfigMap.Data["syscomponents.yaml"]
		err = yaml.Unmarshal([]byte(content), &cachedSysComponentConfigs)
		if err != nil {
			klog.Errorf("failed to unmarshal syscomponents config: %v", err)
			return cachedSysComponentConfigs, err
		}
	}
	return cachedSysComponentConfigs, nil
}
