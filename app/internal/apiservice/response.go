package apiservice

type APIResponse[K any] struct {
	Result string `json:"result"`
	Data   *K     `json:"data,omitempty"`
}

func NewAPIResponse[K any](data K) *APIResponse[K] {
	return &APIResponse[K]{Result: "success", Data: &data}
}
