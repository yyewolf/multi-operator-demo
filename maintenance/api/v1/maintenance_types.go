/*
Copyright 2025.

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

package v1

import (
	"library"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	routev1 "multi.ch/route/api/v1"
)

// MaintenanceSpec defines the desired state of Maintenance.
type MaintenanceSpec struct {
	// +required
	Replaces *MaintenanceTargetReference `json:"replaces"`
}

type MaintenanceTargetReference struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Name       string `json:"name"`
}

// MaintenanceStatus defines the observed state of Maintenance.
type MaintenanceStatus struct {
	routev1.RouteContractInjector `json:",inline"`

	library.Status `json:",inline"`
}

func (maintenance *Maintenance) GetStatus() *library.Status {
	return &maintenance.Status.Status
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Maintenance is the Schema for the maintenances API.
type Maintenance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MaintenanceSpec   `json:"spec,omitempty"`
	Status MaintenanceStatus `json:"status,omitempty"`
}

var _ library.ControllerResource = &Maintenance{}

// +kubebuilder:object:root=true

// MaintenanceList contains a list of Maintenance.
type MaintenanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Maintenance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Maintenance{}, &MaintenanceList{})
}
