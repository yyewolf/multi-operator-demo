package v1

import (
	"library"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AppSpec defines the desired state of App.
type AppSpec struct {
	// +required
	// +kubebuilder:validation:Minimum=1024
	// +kubebuilder:validation:Maximum=65535
	Port int32 `json:"port"`

	// +required
	Command string `json:"command"`
}

// AppStatus defines the observed state of App.
type AppStatus struct {
	library.Status `json:",inline"`
}

func (app *App) GetStatus() *library.Status {
	return &app.Status.Status
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// App is the Schema for the apps API.
type App struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppSpec   `json:"spec,omitempty"`
	Status AppStatus `json:"status,omitempty"`
}

var _ library.ControllerResource = &App{}

// +kubebuilder:object:root=true

// AppList contains a list of App.
type AppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []App `json:"items"`
}

func init() {
	SchemeBuilder.Register(&App{}, &AppList{})
}
