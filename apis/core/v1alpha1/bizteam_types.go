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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type TeamResourceLimitation corev1.ResourceList

// BizTeamSpec defines the desired state of BizTeam
type BizTeamSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ResourceLimitation defines the resource limitation of the BizTeam
	ResourceLimitation TeamResourceLimitation `json:"resourceLimitation,omitempty"`
}

type BizTeamPhase string

const (
	BizTeamActive      BizTeamPhase = "Active"
	BizTeamCreating    BizTeamPhase = "Creating"
	BizTeamTerminating BizTeamPhase = "Terminating"
)

// BizTeamStatus defines the observed state of BizTeam
type BizTeamStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Phase                  BizTeamPhase           `json:"phase"`
	Users                  []string               `json:"users,omitempty"`
	LimitedResource        TeamResourceLimitation `json:"limitedResource,omitempty"`
	RequestedResource      TeamResourceLimitation `json:"requestedResource,omitempty"`
	ResourcePressureStatus string                 `json:"resourcePressureStatus,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
// +genclient
// +genclient:nonNamespaced

// BizTeam is the Schema for the bizteams API
type BizTeam struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BizTeamSpec   `json:"spec,omitempty"`
	Status BizTeamStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BizTeamList contains a list of BizTeam
type BizTeamList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BizTeam `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BizTeam{}, &BizTeamList{})
}
