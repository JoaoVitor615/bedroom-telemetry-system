package repository

import (
	"context"

	"bedroom-api/internal/domain"
	"bedroom-api/internal/domain/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// MongoAdapter implements domain.IDoorRepository
type MongoAdapter struct {
	collection *mongo.Collection
}

// NewMongoAdapter creates a new repository adapter for MongoDB v2
func NewMongoAdapter(db *mongo.Database, collectionName string) domain.IDoorRepository {
	return &MongoAdapter{
		collection: db.Collection(collectionName),
	}
}

func (a *MongoAdapter) SaveEvent(ctx context.Context, event repository.DoorEvent) error {
	// Map Domain Entity to BSON
	doc := bson.M{
		"state":     event.State,
		"timestamp": event.Timestamp,
	}
	_, err := a.collection.InsertOne(ctx, doc)
	return err
}

func (a *MongoAdapter) GetHistory(ctx context.Context, limit int64) ([]repository.DoorEvent, error) {
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}).SetLimit(limit)

	cursor, err := a.collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode into an anonymous struct using v2 BSON primitive
	var docs []struct {
		State     string        `bson:"state"`
		Timestamp bson.DateTime `bson:"timestamp"`
	}

	if err = cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	var events []repository.DoorEvent
	for _, doc := range docs {
		events = append(events, repository.DoorEvent{
			State:     doc.State,
			Timestamp: doc.Timestamp.Time(),
		})
	}

	return events, nil
}
