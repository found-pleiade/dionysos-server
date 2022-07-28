package utils

type CreateResponse struct {
	URI      string `json:"uri"`
	Password string `json:"password,omitempty"`
}
