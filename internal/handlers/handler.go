package handlers

import (
	"net/http"
	"snake-tournament/internal/iface"
	"snake-tournament/models/dto"

	"github.com/alexsuslov/ehttp"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

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
	service iface.IGameService
}

func NewRecordHandler(service iface.IGameService, router *httprouter.Router) GameHandler {
	handler := GameHandler{service: service}

	router.HandlerFunc(http.MethodPost, startUrl,
		ehttp.Middleware(handler.Start))
	router.HandlerFunc(http.MethodGet, enterToGameUrl,
		ehttp.Middleware(handler.EnterToGame))
	router.HandlerFunc(http.MethodGet, findGamesForCurrentUserUrl,
		ehttp.Middleware(handler.FindGamesForCurrentUser))
	router.HandlerFunc(http.MethodGet, findActiveUrl,
		ehttp.Middleware(handler.FindActiveGames))
	router.HandlerFunc(http.MethodPost, recordsCreateUrl,
		ehttp.Middleware(handler.PasteResults))
	router.HandlerFunc(http.MethodGet, checkRenewAllowedUrl,
		ehttp.Middleware(handler.CheckRenewAllowed))
	router.HandlerFunc(http.MethodPost, selectPrizeUrl,
		ehttp.Middleware(handler.SetPrize))

	return handler
}

// Start game
// @Summary Find game by params or create if not found
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param data body dto.GameCreateRequest true "create record struct"
// @Tags snake
// @Success 200 {object} dto.GameDto
// @Failure 500 {object} ehttp.Error
// @Router /api/v1/snake-tournament/start [post]
func (h *GameHandler) Start(response ehttp.ResponseJson, request *ehttp.Request) error {
	defer request.Body.Close()

	var user dto.User
	err := request.TokenUser(&user)
	if err != nil {
		return ehttp.StatusForbidden("Unauthorized", err)
	}

	var requestEntity dto.GameCreateRequest
	err = request.Unmarshal(&requestEntity)

	if err != nil {
		return ehttp.StatusBadRequest("Invalid json", err)
	}

	logrus.
		WithField("players", requestEntity.PlayersAmount).
		WithField("user", user.Email).
		WithField("userId", user.Id).
		Info("start")

	record, err := h.service.Start(request.Context(), requestEntity, user)

	if err != nil {
		return err
	}

	response.WriteHeader(http.StatusOK)
	return response.WriteJson(record)
}

// Enter game by id
// @Summary Enter game by id
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param type path string true "id"
// @Tags snake
// @Success 200 {object} dto.GameDto
// @Failure 500 {object} ehttp.Error
// @Router /api/v1/snake-tournament/start/{id} [post]
func (h *GameHandler) EnterToGame(response ehttp.ResponseJson, request *ehttp.Request) error {
	defer request.Body.Close()

	params := request.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	id := params.ByName("id")

	var user dto.User
	err := request.TokenUser(&user)
	if err != nil {
		return ehttp.StatusForbidden("Unauthorized", err)
	}

	logrus.
		WithField("gameID", id).
		WithField("user", user.Email).
		WithField("userId", user.Id).
		Info("EnterToGame")

	record, err := h.service.EnterToGame(request.Context(), id, user)

	if err != nil {
		return err
	}

	response.WriteHeader(http.StatusOK)
	return response.WriteJson(record)
}

// Get games by user
// @Summary Get own games
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Tags snake
// @Success 200 {object} []dto.GameDto
// @Failure 500 {object} ehttp.Error
// @Router /api/v1/snake-tournament/find-my [get]
func (h *GameHandler) FindGamesForCurrentUser(response ehttp.ResponseJson, request *ehttp.Request) error {
	defer request.Body.Close()

	var user dto.User
	err := request.TokenUser(&user)
	if err != nil {
		return ehttp.StatusForbidden("Unauthorized", err)
	}

	logrus.
		WithField("user", user.Email).
		WithField("userId", user.Id).
		Info("FindGamesForCurrentUser")

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
// @Tags snake
// @Success 200 {object} []dto.GameDto
// @Failure 400 {object} ehttp.Error
// @Router /api/v1/snake-tournament/find-active [get]
func (h *GameHandler) FindActiveGames(response ehttp.ResponseJson, request *ehttp.Request) error {
	defer request.Body.Close()
	var user dto.User
	err := request.TokenUser(&user)
	if err != nil {
		return ehttp.StatusForbidden("Unauthorized", err)

	}

	logrus.
		WithField("user", user.Email).
		WithField("userId", user.Id).
		Info("FindActiveGames")

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
// @Param type path string true "id"
// @Tags snake
// @Success 200 {object} dto.GameDto
// @Failure 400 {object} ehttp.Error
// @Router /api/v1/snake-tournament/create-record/:id [post]
func (h *GameHandler) PasteResults(response ehttp.ResponseJson, request *ehttp.Request) error {
	defer request.Body.Close()

	var user dto.User
	err := request.TokenUser(&user)
	if err != nil {
		return ehttp.StatusForbidden("Unauthorized", err)
	}

	var requestEntity dto.RecordCreateRequest
	err = request.Unmarshal(&requestEntity)
	if err != nil {
		return ehttp.StatusBadRequest("Invalid json", err)
	}

	params := request.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	gameId := params.ByName("id")

	logrus.
		WithField("user", user.Email).
		WithField("userId", user.Id).
		WithField("gameId", gameId).
		Info("PasteResults")

	records, err := h.service.PasteResults(request.Context(), gameId, requestEntity, user)
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
// @Param type path string true "id"
// @Tags snake
// @Success 200 {object} dto.GameDto
// @Failure 500 {object} ehttp.Error
// @Router /api/v1/snake-tournament/check-renew-record-allowed/:id [post]
func (h *GameHandler) CheckRenewAllowed(response ehttp.ResponseJson, request *ehttp.Request) error {
	defer request.Body.Close()

	var user dto.User
	err := request.TokenUser(&user)

	if err != nil {
		return ehttp.StatusForbidden("Unauthorized", err)
	}

	params := request.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	gameId := params.ByName("id")

	logrus.
		WithField("user", user.Email).
		WithField("userId", user.Id).
		WithField("gameId", gameId).
		Info("CheckRenewAllowed")

	result, err := h.service.CheckAvailableRenew(request.Context(), gameId, user)

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
// @Param type path string true "id"
// @Tags snake
// @Success 200 {object} dto.GameDto
// @Failure 401 {object} ehttp.Error
// @Failure 500 {object} ehttp.Error
// @Router /api/v1/snake-tournament/select-prize/:id [post]
func (h *GameHandler) SetPrize(response ehttp.ResponseJson, request *ehttp.Request) error {
	defer request.Body.Close()
	var user dto.User

	if err := request.TokenUser(&user); err != nil {
		return ehttp.StatusForbidden("Unauthorized", err)
	}

	var requestEntity dto.SelectPrizeDto
	if err := request.Unmarshal(&requestEntity); err != nil {
		return ehttp.StatusBadRequest("Invalid json", err)
	}

	params := request.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	gameId := params.ByName("id")

	logrus.
		WithField("user", user.Email).
		WithField("userId", user.Id).
		WithField("gameId", gameId).
		Info("SetPrize")

	record, err := h.service.SetPrize(request.Context(), gameId, requestEntity, user)
	if err != nil {
		return err
	}

	response.WriteHeader(http.StatusOK)
	return response.WriteJson(record)
}
