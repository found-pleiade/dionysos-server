package utils

import "github.com/Brawdunoir/dionysos-server/constants"

type URIResponse struct {
	URI string `json:"uri"`
}

func CreateURIResponse(uri string) *URIResponse {
	return &URIResponse{
		URI: constants.BasePath + uri,
	}
}
