package repository

import (
	"context"
	"fmt"
	"net/http"
	"snake-tournament/internal/iface"

	"github.com/alexsuslov/ehttp"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	collection *mongo.Collection
}

func New(database *mongo.Database, collection string) *Repository {
	return &Repository{
		collection: database.Collection(collection),
	}
}

func (d *Repository) Create(ctx context.Context, model iface.IModel) error {
	result, err := d.collection.InsertOne(ctx, model)
	if err != nil {
		return fmt.Errorf("failed to create record due to: %v", err)
	}

	oid, ok := result.InsertedID.(bson.ObjectID)
	if ok {
		model.SetId(oid.Hex())
		return nil
	}

	return fmt.Errorf("failed to convert objectId to hex. probable oid: %s", oid)
}

func (d *Repository) FindAll(ctx context.Context, filter interface{}, models any) (err error) {
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

func (d *Repository) Find(ctx context.Context, filter interface{}, model any) (err error) {
	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return result.Err()
	}

	if err = result.Decode(model); err != nil {
		return fmt.Errorf("failed to decode record from DB due to error: %v", err)
	}

	return nil
}

func (d *Repository) FindById(ctx context.Context, id string, res any) (err error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert hex to objectID, hex: %s", id)
	}
	filter := bson.M{"_id": oid}
	result := d.collection.FindOne(ctx, filter)

	if result.Err() != nil {
		return &ehttp.Error{
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

func (d *Repository) Update(ctx context.Context, res iface.IModel) error {
	oid, err := bson.ObjectIDFromHex(res.GetId())
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
		return &ehttp.Error{
			Message: "Not found",
			Code:    http.StatusNotFound,
			Err:     err,
		}
	}

	return nil
}

func (d *Repository) Delete(ctx context.Context, objectID bson.ObjectID) error {
	//objectID, err := bson.ObjectIDFromHex(id)
	//if err != nil {
	//	return fmt.Errorf("failed to convert user ID to ObjectID. ID=%v", id)
	//}

	result, err := d.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to execute update user query due to: %v", err)
	}

	if result.DeletedCount == 0 {
		return &ehttp.Error{
			Message: "Not found",
			Code:    http.StatusNotFound,
			Err:     err,
		}
	}

	return nil
}

func (d *Repository) DeleteAll(ctx context.Context, ids []string) error {
	objectIds := make([]bson.ObjectID, len(ids))

	var err error

	for i := range ids {
		objectIds[i], err = bson.ObjectIDFromHex(ids[i])

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
