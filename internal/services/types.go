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
	"github.com/wutong-paas/region/apis/core/v1alpha1"
	"github.com/wutong-paas/region/generated/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	controllerruntime "sigs.k8s.io/controller-runtime"
)

var defaultServicer *CommonServicer

type CommonServicer struct {
	KubeConfig   *rest.Config
	KubeClient   kubernetes.Interface
	RegionClient versioned.Interface
}

func DefaultServicer() *CommonServicer {
	if defaultServicer == nil {
		kubeConfig := controllerruntime.GetConfigOrDie()
		defaultServicer = &CommonServicer{
			KubeConfig:   kubeConfig,
			KubeClient:   kubernetes.NewForConfigOrDie(kubeConfig),
			RegionClient: versioned.NewForConfigOrDie(kubeConfig),
		}
	}
	return defaultServicer
}

// syscomponent

type SysComponentStatus string

const (
	SysComponentInstalled    SysComponentStatus = "Installed"
	SysComponentUnInstalled  SysComponentStatus = "Uninstalled"
	SysComponentInstalling   SysComponentStatus = "Installing"
	SysComponentUnInstalling SysComponentStatus = "Uninstalling"
	SysComponentUpgrading    SysComponentStatus = "Upgrading"
	// SysComponentUnAbnormal   SysComponentStatus = "Abnormal"
	// SysComponentUnUnknown    SysComponentStatus = "Unknown"
)

type SysComponent struct {
	Name              string             `json:"name"`
	Namespace         string             `json:"namespace"`
	Description       string             `json:"description"`
	AvailableVersions []string           `json:"availableVersions"`
	CurrentVersion    string             `json:"currentVersion"`
	Status            SysComponentStatus `json:"status"`
}

type SysComponentConfig struct {
	// Name              string                                      `json:"name,omitempty"`
	Namespace         string                                       `json:"namespace,omitempty"`
	Description       string                                       `json:"description,omitempty"`
	InstallWay        string                                       `json:"installWay,omitempty"`
	AvailableVersions map[string]*v1alpha1.SysComponentVersionInfo `json:"availableVersions,omitempty"`
}

// bizcomponnet
