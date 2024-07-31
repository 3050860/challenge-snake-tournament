package ehttp

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ExtendedResponse struct {
	http.ResponseWriter
}

func (r ExtendedResponse) WriteJson(responseModel any) error {
	r.Header().Set("Content-Type", "application/json")
	userBytes, err := json.Marshal(responseModel)
	if err != nil {
		return fmt.Errorf("failed to marshall data. error: %w", err)
	}

	r.Write(userBytes)

	return nil
}
