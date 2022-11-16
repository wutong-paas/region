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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	// Install the SysComponent by helm install
	InstallHelm = "helm"

	// Install the SysComponent by client apply
	InstallApply = "apply"
)

type SysComponentVersionInfo struct {
	HelmRepoName  string `json:"helmRepoName,omitempty"`
	HelmChartName string `json:"helmChartName,omitempty"`
	HelmRepoUrl   string `json:"helmRepoUrl,omitempty"`

	ApplyFileUrl     string `json:"applyFileUrl,omitempty"`
	ApplyFileContent string `json:"applyFileContent,omitempty"`
}

// SysComponentSpec defines the desired state of SysComponent
type SysComponentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Description string `json:"description,omitempty"`

	AvailableVersions map[string]*SysComponentVersionInfo `json:"availableVersions,omitempty"`
	InstallWay        string                              `json:"installWay,omitempty"`
	CurrentVersion    string                              `json:"currentVersion,omitempty"`
}

const (
	SysComponentInstalled      = "Installed"
	SysComponentUnInstalled    = "UnInstalled"
	SysComponentInstalling     = "Installing"
	SysComponentPendingInstall = "PendingInstall"
	SysComponentUnInstalling   = "UnInstalling"
	SysComponentUpgrading      = "Upgrading"
	SysComponentPendingUpgrade = "PendingUpgrade"
	// SysComponentAbnormal       = "Abnormal"
	// SysComponentUnknown        = "Unknown"
)

// SysComponentStatus defines the observed state of SysComponent
type SysComponentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Phase
	Phase string `json:"phase,omitempty"`

	// Message is a message for SysComponent installing or uninstalling
	Message string `json:"message,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +genclient

// SysComponent is the Schema for the syscomponents API
type SysComponent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SysComponentSpec   `json:"spec,omitempty"`
	Status SysComponentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SysComponentList contains a list of SysComponent
type SysComponentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SysComponent `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SysComponent{}, &SysComponentList{})
}
