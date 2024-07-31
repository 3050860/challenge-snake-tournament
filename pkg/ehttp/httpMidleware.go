package ehttp

import (
	"errors"
	"log"
	"net/http"
)

type CustomError interface {
	Marshal() []byte
	Status() int
}

type UserService interface {
	Extract(r ExtendedRequest, user any) error
}

type extendedAppHandler func(ExtendedResponse, ExtendedRequest) error

func ExtendedMiddleware(h extendedAppHandler, s UserService) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		var appErr CustomError
		newResponse := ExtendedResponse{
			ResponseWriter: response,
		}

		newRequest := ExtendedRequest{
			Request:             *request,
			extractUserFunction: s.Extract,
		}

		err := h(newResponse, newRequest)
		if err != nil {
			if errors.As(err, &appErr) {
				newResponse.WriteHeader(appErr.Status())
				newResponse.Write(appErr.Marshal())
				log.Print(string(appErr.Marshal()))
				return
			}
			newResponse.WriteHeader(http.StatusInternalServerError)
		}
	}
}
