package repository

import (
	"context"
	"snake-tournament/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// GamesDatabase struct represents the database operations for game data
type GamesDatabase struct {
	Repository
}

// NewGames creates a new instance of GamesDatabase
func NewGames(database *mongo.Database, collection string) *GamesDatabase {
	return &GamesDatabase{
		Repository: *New(database, collection),
	}
}

// FindAvailableToEnterGame finds a game that is available for a user to join
// It looks for games with exact player count and where the user hasn't joined yet
func (d *GamesDatabase) FindAvailableToEnterGame(ctx context.Context, playerCount int, userId string) *models.Game {
	var game models.Game
	now := time.Now()
	filter := bson.D{
		{Key: "players_amount", Value: playerCount},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "close_time", Value: bson.D{{Key: "$gt", Value: now}}}},
			bson.D{{Key: "close_time", Value: nil}},
		}},
		{Key: "$expr", Value: bson.D{
			{Key: "$lt", Value: bson.A{bson.D{{Key: "$size", Value: "$records"}}, playerCount}},
		}},
		{Key: "records", Value: bson.D{{Key: "$not", Value: bson.D{{Key: "$elemMatch", Value: bson.D{{Key: "user_id", Value: userId}}}}}}},
	}
	opts := options.FindOne().SetSort(bson.D{{Key: "_id", Value: 1}})
	if err := d.Find(ctx, filter, &game, opts); err != nil {
		return nil
	}
	return &game
}

// FindAvailableToEnterGames finds all games that are available for a user to join
// It looks for games where the user hasn't joined yet and there are still spots available
func (d *GamesDatabase) FindAvailableToEnterGames(ctx context.Context, userId string) []models.Game {
	var games []models.Game
	now := time.Now()

	filter := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "close_time", Value: bson.D{{Key: "$gt", Value: now}}}},
			bson.D{{Key: "close_time", Value: nil}},
		}},
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

// FindGameByIncludeUser finds all games that include a specific user
// This is used to get all games that a user has participated in
func (d *GamesDatabase) FindGameByIncludeUser(ctx context.Context, userId string) ([]models.Game, error) {
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
