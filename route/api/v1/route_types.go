package v1

import (
	"library"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// RouteSpec defines the desired state of Route.
type RouteSpec struct {
	// +required
	// +kubebuilder:validation:MinItems=1
	Hostnames []gatewayv1.PreciseHostname `json:"hostnames,omitempty"`

	// +required
	// +kubebuilder:validation:MinItems=1
	TargetRefs []*RouteTargetReference `json:"targetRefs,omitempty"`
}

type RouteTargetReference struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Name       string `json:"name"`
}

// RouteStatus defines the observed state of Route.
type RouteStatus struct {
	library.Status `json:",inline"`
}

func (route *Route) GetStatus() *library.Status {
	return &route.Status.Status
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Route is the Schema for the routes API.
type Route struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RouteSpec   `json:"spec,omitempty"`
	Status RouteStatus `json:"status,omitempty"`
}

var _ library.ControllerResource = &Route{}

// +kubebuilder:object:root=true

// RouteList contains a list of Route.
type RouteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Route `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Route{}, &RouteList{})
}
