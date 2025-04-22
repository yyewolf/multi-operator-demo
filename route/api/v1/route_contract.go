package v1

type RouteContractInjector struct {
	// +optional
	Contract RouteContract `json:"routeContract"`
}

type RouteContract struct {
	// +required
	ServiceRef RouteContractLocalServiceRef `json:"serviceRef"`
}

type RouteContractLocalServiceRef struct {
	// +required
	Name string `json:"name"`
	// +required
	Port int `json:"port"`
}
