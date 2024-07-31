package ehttp

import (
	"encoding/json"
)

type HttpError struct {
	Err     error  `json:"-"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"-"`
}

func (e *HttpError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}

	return e.Message
}

func (e *HttpError) Marshal() []byte {
	bytes, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return bytes
}

func (e *HttpError) Status() int {
	return e.Code
}
