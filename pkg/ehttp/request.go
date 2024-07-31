package ehttp

import (
	"encoding/json"
	"net/http"
)

type ExtendedRequest struct {
	http.Request
	extractUserFunction func(ExtendedRequest, any) error
}

func (r ExtendedRequest) ExtractDto(requestModel any) error {
	err := json.NewDecoder(r.Body).Decode(requestModel)
	if err != nil {
		return &HttpError{
			Message: "Bad request",
			Code:    http.StatusBadRequest,
		}
	}

	return nil
}

func (r ExtendedRequest) ExtractUser(user any) error {
	return r.extractUserFunction(r, user)
}
