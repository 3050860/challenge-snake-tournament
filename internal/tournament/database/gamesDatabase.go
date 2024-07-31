package database

import (
	"context"
	"snake-tournament/models"
	"snake-tournament/pkg/mongodb_storage"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GamesDatabase struct {
	mongodb_storage.BaseStorage
}

func NewGamesDatabase(database *mongo.Database, collection string) *GamesDatabase {
	return &GamesDatabase{
		BaseStorage: *mongodb_storage.NewStorage(database, collection),
	}
}

func (db GamesDatabase) FindAvailableToEnterGame(ctx context.Context, playerCount int, userId string) *models.Game {
	var game models.Game
	filter := bson.D{
		{Key: "players_amount", Value: playerCount},
		{Key: "$expr", Value: bson.D{
			{Key: "$lt", Value: bson.A{bson.D{{Key: "$size", Value: "$records"}}, playerCount}},
		}},
		{Key: "records", Value: bson.D{{Key: "$not", Value: bson.D{{Key: "$elemMatch", Value: bson.D{{Key: "user_id", Value: userId}}}}}}},
	}
	if err := db.Find(ctx, filter, &game); err != nil {
		return nil
	}
	return &game
}

func (d GamesDatabase) FindAvailableToEnterGames(ctx context.Context, userId string) []models.Game {
	var games []models.Game

	filter := bson.D{
		{Key: "$expr", Value: bson.D{
			{Key: "$lt", Value: bson.A{bson.D{{Key: "$size", Value: "$records"}}, "$players_amount"}},
		}},
		{Key: "records", Value: bson.D{{Key: "$not", Value: bson.D{{Key: "$elemMatch", Value: bson.D{{Key: "user_id", Value: userId}}}}}}},
	}

	if err := d.FindAll(ctx, filter, &games); err != nil {
		return make([]models.Game, 0)
	}

	return games
}

func (d GamesDatabase) FindGameByIncludeUser(ctx context.Context, userId string) ([]models.Game, error) {
	filter := bson.M{
		"records": bson.M{
			"$elemMatch": bson.M{
				"user_id": userId,
			},
		},
	}

	var games []models.Game

	if err := d.FindAll(ctx, filter, &games); err != nil {
		return nil, err
	}

	return games, nil
}
