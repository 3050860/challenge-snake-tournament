package iface

import (
	"context"
	"snake-tournament/models"
	"snake-tournament/models/dto"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type IRepository interface {
	FindById(ctx context.Context, id string, res any) (err error)
	Update(ctx context.Context, res IModel) error
	Create(ctx context.Context, model IModel) error
	Delete(ctx context.Context, objectID bson.ObjectID) error
}

type IGameRepository interface {
	IRepository

	FindAvailableToEnterGame(ctx context.Context, playerCount int, userId string) *models.Game
	FindGameByIncludeUser(ctx context.Context, userId string) ([]models.Game, error)
	FindAvailableToEnterGames(ctx context.Context, userId string) []models.Game
}

type IGameService interface {
	Start(ctx context.Context, dto dto.GameCreateRequest, user dto.User) (*dto.GameDto, error)
	EnterToGame(ctx context.Context, id string, user dto.User) (*dto.GameDto, error)
	GetGamesForCurrentUser(ctx context.Context, user dto.User) ([]dto.ResultGameDto, error)
	GetActiveGame(ctx context.Context, user dto.User) ([]dto.GameDto, error)
	PasteResults(ctx context.Context, id string, dto dto.RecordCreateRequest, user dto.User) ([]dto.RecordDto, error)
	CheckAvailableRenew(ctx context.Context, id string, user dto.User) (bool, error)
	SetPrize(ctx context.Context, id string, request dto.SelectPrizeDto, user dto.User) (*dto.ResultGameDto, error)
}
