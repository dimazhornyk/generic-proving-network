package common

type ProvingMessage struct {
	RequestID string `json:"request_id"`
	Data      []byte `json:"data"`
}

type ProvingResponse struct {
	Proof []byte `json:"proof"`
}

type ValidationMessage struct {
	RequestID string `json:"request_id"`
	Proof     []byte `json:"proof"`
	Data      []byte `json:"data"`
}

type ValidationResponse struct {
	Valid bool `json:"valid"`
}
