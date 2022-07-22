package utils

type URIResponse struct {
	URI string `json:"uri"`
}

func CreateURIResponse(uri string) *URIResponse {
	return &URIResponse{
		URI: uri,
	}
}
