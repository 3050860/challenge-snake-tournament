package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewClient(ctx context.Context, url, database string, authSource string, username string, password string) (db *mongo.Database, err error) {
	credential := options.Credential{
		AuthSource: authSource,
		Username:   username,
		Password:   password,
	}

	clientOptions := options.Client().ApplyURI(url).SetAuth(credential)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongoDB due to error: %v", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping mongoDB due to error: %v", err)
	}

	return client.Database(database), nil
}
