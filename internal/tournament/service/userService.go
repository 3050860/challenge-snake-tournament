package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"snake-tournament/internal/config"
	"snake-tournament/pkg/ehttp"
)

type UserService struct {
	config *config.Config
}

func NewUserService(config *config.Config) *UserService {
	return &UserService{
		config: config,
	}
}

func (s UserService) Extract(request ehttp.ExtendedRequest, user any) error {
	authorizationToken := request.Header.Get("Authorization")

	if authorizationToken == "" {
		return &ehttp.HttpError{
			Message: "Not authorized",
			Code:    http.StatusForbidden,
			Err:     nil,
		}
	}

	requestURL := fmt.Sprintf("http://%s/api/v1/auth/current", s.config.UserService.Host)
	client := http.Client{}
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil
	}

	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {authorizationToken},
	}

	res, err := client.Do(req)
	if err != nil {
		return nil
	}

	err = json.NewDecoder(res.Body).Decode(user)
	if err != nil {
		return nil
	}

	return nil
}
