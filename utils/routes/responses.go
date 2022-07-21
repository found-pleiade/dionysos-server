package utils

import "github.com/Brawdunoir/dionysos-server/variables"

type URIResponse struct {
	URI string `json:"uri"`
}

func CreateURIResponse(uri string) *URIResponse {
	return &URIResponse{
		URI: variables.BasePath + uri,
	}
}
