package mongodb_storage

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"snake-tournament/models"
)

type Model interface {
	GetId() string
	SetId(string)
}

type BaseStorage struct {
	collection *mongo.Collection
}

func NewStorage(database *mongo.Database, collection string) *BaseStorage {
	return &BaseStorage{
		collection: database.Collection(collection),
	}
}

func (d *BaseStorage) Create(ctx context.Context, model Model) error {
	result, err := d.collection.InsertOne(ctx, model)
	if err != nil {
		return fmt.Errorf("failed to create record due to: %v", err)
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		model.SetId(oid.Hex())
		return nil
	}

	return fmt.Errorf("failed to convert objectId to hex. probable oid: %s", oid)
}

func (d *BaseStorage) FindAll(ctx context.Context, filter interface{}, models any) (err error) {
	cursor, err := d.collection.Find(ctx, filter)
	if err != nil {
		return err
	}
	if cursor.Err() != nil {
		return fmt.Errorf("failed to find all users due to: %v", cursor.Err())
	}
	if err := cursor.All(ctx, models); err != nil {
		return fmt.Errorf("failed to read all documents from cursor")
	}
	return nil
}

func (d *BaseStorage) Find(ctx context.Context, filter interface{}, model any) (err error) {
	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return result.Err()
	}

	if err = result.Decode(model); err != nil {
		return fmt.Errorf("failed to decode record from DB due to error: %v", err)
	}

	return nil
}

func (d *BaseStorage) FindById(ctx context.Context, id string, res any) (err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert hex to objectID, hex: %s", id)
	}
	filter := bson.M{"_id": oid}
	result := d.collection.FindOne(ctx, filter)

	if result.Err() != nil {
		return &models.Error{
			Message: "Not found",
			Code:    http.StatusNotFound,
			Err:     err,
		}
	}

	if err = result.Decode(res); err != nil {
		return fmt.Errorf("failed to decode record(id:%s) from DB due to error: %v", id, err)
	}
	return nil
}

func (d *BaseStorage) Update(ctx context.Context, res Model) error {
	oid, err := primitive.ObjectIDFromHex(res.GetId())
	if err != nil {
		return err
	}
	filter := bson.M{"_id": oid}

	res.SetId("")

	result, err := d.collection.ReplaceOne(ctx, filter, res)

	res.SetId(oid.Hex())

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return &models.Error{
			Message: "Not found",
			Code:    http.StatusNotFound,
			Err:     err,
		}
	}

	return nil
}

func (d *BaseStorage) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert user ID to ObjectID. ID=%v", id)
	}

	result, err := d.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to execute update user query due to: %v", err)
	}

	if result.DeletedCount == 0 {
		return &models.Error{
			Message: "Not found",
			Code:    http.StatusNotFound,
			Err:     err,
		}
	}

	return nil
}

func (d *BaseStorage) DeleteAll(ctx context.Context, ids []string) error {
	objectIds := make([]primitive.ObjectID, len(ids))

	var err error

	for i := range ids {
		objectIds[i], err = primitive.ObjectIDFromHex(ids[i])

		if err != nil {
			return err
		}
	}

	_, err = d.collection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": objectIds}})

	if err != nil {
		return fmt.Errorf("failed to execute update user query due to: %v", err)
	}

	return nil
}
