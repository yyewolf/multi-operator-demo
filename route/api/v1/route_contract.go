package v1

import (
	"reflect"
)

type RouteContractInjector struct {
	// +optional
	RouteContract RouteContract `json:"routeContract"`
}

func (contract *RouteContractInjector) Get() RouteContract {
	return contract.RouteContract
}

func (contract *RouteContractInjector) Set(new RouteContract) bool {
	changed := !reflect.DeepEqual(new, contract.RouteContract)
	if changed {
		contract.RouteContract = new
	}

	return changed
}

type RouteContract struct {
	// +optional
	ServiceRef *RouteContractLocalServiceRef `json:"serviceRef,omitempty"`

	// +optional
	BackendRef *RouteContractLocalBackendRef `json:"backendRef,omitempty"`
}

type RouteContractLocalServiceRef struct {
	// +required
	Name string `json:"name"`
	// +required
	Port int `json:"port"`
}

type RouteContractLocalBackendRef struct {
	// +required
	Name string `json:"name"`
	// +required
	Port int `json:"port"`
}
