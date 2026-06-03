package handler

import (
	"context"
	"net/http"
	"snake-tournament/models"
	"snake-tournament/models/dto"
	"snake-tournament/pkg/ehttp"
	"snake-tournament/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

type GameService interface {
	Start(ctx context.Context, dto dto.GameCreateRequest, user dto.User) (*dto.GameDto, error)
	EnterToGame(ctx context.Context, id string, user dto.User) (*dto.GameDto, error)
	GetGamesForCurrentUser(ctx context.Context, user dto.User) ([]dto.ResultGameDto, error)
	GetActiveGame(ctx context.Context, user dto.User) ([]dto.GameDto, error)
	PasteResults(ctx context.Context, id string, dto dto.RecordCreateRequest, user dto.User) ([]dto.RecordDto, error)
	SetPrize(ctx context.Context, id string, request dto.SelectPrizeDto, user dto.User) (*dto.ResultGameDto, error)
	CheckAvailableRenew(ctx context.Context, id string, user dto.User) (bool, error)
}

var (
	startUrl                   = "/api/v1/snake-tournament/start"
	enterToGameUrl             = "/api/v1/snake-tournament/start/:id"
	findGamesForCurrentUserUrl = "/api/v1/snake-tournament/find-my"
	findActiveUrl              = "/api/v1/snake-tournament/find-active"
	recordsCreateUrl           = "/api/v1/snake-tournament/create-record/:id"
	checkRenewAllowedUrl       = "/api/v1/snake-tournament/check-renew-record-allowed/:id"
	selectPrizeUrl             = "/api/v1/snake-tournament/select-prize/:id"
)

type GameHandler struct {
	service GameService
}

func NewRecordHandler(service GameService, router *httprouter.Router, userService ehttp.UserService) GameHandler {
	handler := GameHandler{service: service}

	router.HandlerFunc(http.MethodPost, startUrl, ehttp.ExtendedMiddleware(handler.Start, userService))
	router.HandlerFunc(http.MethodGet, enterToGameUrl, ehttp.ExtendedMiddleware(handler.EnterToGame, userService))
	router.HandlerFunc(http.MethodGet, findGamesForCurrentUserUrl, ehttp.ExtendedMiddleware(handler.FindGamesForCurrentUser, userService))
	router.HandlerFunc(http.MethodGet, findActiveUrl, ehttp.ExtendedMiddleware(handler.FindActiveGames, userService))
	router.HandlerFunc(http.MethodPost, recordsCreateUrl, ehttp.ExtendedMiddleware(handler.PasteResults, userService))
	router.HandlerFunc(http.MethodGet, checkRenewAllowedUrl, ehttp.ExtendedMiddleware(handler.CheckRenewAllowed, userService))
	router.HandlerFunc(http.MethodPost, selectPrizeUrl, ehttp.ExtendedMiddleware(handler.SetPrize, userService))

	return handler
}

// Start game
// @Summary Find game by params or create if not found
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param data body dto.GameCreateRequest true "create record struct"
// @Tags Records
// @Success 200 {object} dto.GameDto
// @Failure 400 {object} custom_error.Error
// @Router /api/v1/snake-tournament/start [post]
func (h *GameHandler) Start(response ehttp.ExtendedResponse, request ehttp.ExtendedRequest) error {
	defer request.Body.Close()
	logging.GetLogger().Debugf("Receive %s request on %s url: %s", request.Method, request.RequestURI, request.Body)

	var requestEntity dto.GameCreateRequest
	err := request.ExtractDto(&requestEntity)

	if err != nil {
		return &models.Error{
			Message: "Invalid json",
			Code:    http.StatusBadRequest,
			Err:     nil,
		}
	}

	var user dto.User
	err = request.ExtractUser(&user)

	if err != nil {
		return &models.Error{
			Message: "Unauthorized",
			Code:    http.StatusForbidden,
			Err:     nil,
		}
	}

	logging.GetLogger().Debugf("User %s", user.Email)

	record, err := h.service.Start(request.Context(), requestEntity, user)

	if err != nil {
		return err
	}

	response.WriteHeader(http.StatusOK)
	return response.WriteJson(record)
}

// Enter game
// @Summary Find game by params or create if not found
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param data body dto.GameCreateRequest true "create record struct"
// @Tags Records
// @Success 200 {object} dto.GameDto
// @Failure 400 {object} custom_error.Error
// @Router /api/v1/snake-tournament/start/{id} [post]
func (h *GameHandler) EnterToGame(response ehttp.ExtendedResponse, request ehttp.ExtendedRequest) error {
	defer request.Body.Close()

	logging.GetLogger().Debugf("Receive %s request on %s url: %s", request.Method, request.RequestURI, request.Body)

	params := request.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	id := params.ByName("id")

	var user dto.User
	err := request.ExtractUser(&user)

	if err != nil {
		return &models.Error{
			Message: "Unauthorized",
			Code:    http.StatusForbidden,
			Err:     nil,
		}
	}

	record, err := h.service.EnterToGame(request.Context(), id, user)

	if err != nil {
		return err
	}

	response.WriteHeader(http.StatusOK)
	return response.WriteJson(record)
}

// Get games by user
// @Summary Get games by user
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Tags Records
// @Success 200 {object} []dto.GameDto
// @Failure 400 {object} custom_error.Error
// @Router /api/v1/snake-tournament/find-my [get]
func (h *GameHandler) FindGamesForCurrentUser(response ehttp.ExtendedResponse, request ehttp.ExtendedRequest) error {
	defer request.Body.Close()

	logging.GetLogger().Debugf("Receive %s request on %s url: %s", request.Method, request.RequestURI, request.Body)

	var user dto.User
	err := request.ExtractUser(&user)

	if err != nil {
		return &models.Error{
			Message: "Unauthorized",
			Code:    http.StatusForbidden,
			Err:     err,
		}
	}

	dtos, err := h.service.GetGamesForCurrentUser(request.Context(), user)
	if err != nil {
		return err
	}

	response.WriteHeader(http.StatusOK)
	return response.WriteJson(dtos)
}

// Get active games by user
// @Summary Get active games by user
// @Accept json
// @Produce json
// @Tags Records
// @Success 200 {object} []dto.GameDto
// @Failure 400 {object} custom_error.Error
// @Router /api/v1/snake-tournament/find-active [get]
func (h *GameHandler) FindActiveGames(response ehttp.ExtendedResponse, request ehttp.ExtendedRequest) error {
	defer request.Body.Close()

	// logging.GetLogger().Debugf("Receive %s request on %s url: %s", request.Method, request.RequestURI, request.Body)

	var user dto.User
	err := request.ExtractUser(&user)

	if err != nil {
		return &models.Error{
			Message: "Unauthorized",
			Code:    http.StatusForbidden,
			Err:     err,
		}
	}

	dtos, err := h.service.GetActiveGame(request.Context(), user)
	if err != nil {
		return err
	}

	response.WriteHeader(http.StatusOK)
	return response.WriteJson(dtos)
}

// Send results for end game
// @Summary Send results for end game
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param data body dto.GameCreateRequest true "create record struct"
// @Tags Records
// @Success 200 {object} dto.GameDto
// @Failure 400 {object} custom_error.Error
// @Router /api/v1/snake-tournament/start [post]
func (h *GameHandler) PasteResults(response ehttp.ExtendedResponse, request ehttp.ExtendedRequest) error {
	defer request.Body.Close()

	logging.GetLogger().Debugf("Receive %s request on %s url", request.Method, request.RequestURI)

	var requestEntity dto.RecordCreateRequest
	err := request.ExtractDto(&requestEntity)

	if err != nil {
		return &models.Error{
			Message: "Invalid json",
			Code:    http.StatusBadRequest,
			Err:     nil,
		}
	}

	var user dto.User
	err = request.ExtractUser(&user)

	if err != nil {
		return &models.Error{
			Message: "Unauthorized",
			Code:    http.StatusForbidden,
			Err:     nil,
		}
	}

	params := request.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	id := params.ByName("id")

	records, err := h.service.PasteResults(request.Context(), id, requestEntity, user)

	if err != nil {
		return err
	}

	response.WriteHeader(http.StatusOK)
	return response.WriteJson(records)
}

// Send results for end game
// @Summary Send results for end game
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param data body dto.GameCreateRequest true "create record struct"
// @Tags Records
// @Success 200 {object} dto.GameDto
// @Failure 400 {object} custom_error.Error
// @Router /api/v1/snake-tournament/start [post]
func (h *GameHandler) CheckRenewAllowed(response ehttp.ExtendedResponse, request ehttp.ExtendedRequest) error {
	defer request.Body.Close()

	logging.GetLogger().Debugf("Receive %s request on %s url: %s", request.Method, request.RequestURI, request.Body)

	var user dto.User
	err := request.ExtractUser(&user)

	if err != nil {
		return &models.Error{
			Message: "Unauthorized",
			Code:    http.StatusForbidden,
			Err:     nil,
		}
	}

	params := request.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	id := params.ByName("id")

	result, err := h.service.CheckAvailableRenew(request.Context(), id, user)

	if err != nil {
		return err
	}

	response.WriteHeader(http.StatusOK)
	return response.WriteJson(result)
}

// Send results for end game
// @Summary Send results for end game
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param data body dto.GameCreateRequest true "create record struct"
// @Tags Records
// @Success 200 {object} dto.GameDto
// @Failure 400 {object} custom_error.Error
// @Router /api/v1/snake-tournament/start [post]
func (h *GameHandler) SetPrize(response ehttp.ExtendedResponse, request ehttp.ExtendedRequest) error {
	defer request.Body.Close()

	logging.GetLogger().Debugf("Receive %s request on %s url: %s", request.Method, request.RequestURI, request.Body)

	var requestEntity dto.SelectPrizeDto
	err := request.ExtractDto(&requestEntity)

	if err != nil {
		return &models.Error{
			Message: "Invalid json",
			Code:    http.StatusBadRequest,
			Err:     nil,
		}
	}

	var user dto.User
	err = request.ExtractUser(&user)

	if err != nil {
		return &models.Error{
			Message: "Unauthorized",
			Code:    http.StatusForbidden,
			Err:     nil,
		}
	}

	params := request.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	id := params.ByName("id")

	record, err := h.service.SetPrize(request.Context(), id, requestEntity, user)

	if err != nil {
		return err
	}

	response.WriteHeader(http.StatusOK)
	return response.WriteJson(record)
}
