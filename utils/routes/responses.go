package utils

import "github.com/Brawdunoir/dionysos-server/constants"

type ErrorResponse struct {
	Error string `json:"error"`
}

type UriResponse struct {
	URI string `json:"uri" example:"/api/v0/users/1"`
}

func CreateUriResponse(uri string) *UriResponse {
	return &UriResponse{
		URI: constants.BasePath + uri,
	}
}

func CreateErrorResponse(error string) *ErrorResponse {
	return &ErrorResponse{
		Error: error,
	}
}
