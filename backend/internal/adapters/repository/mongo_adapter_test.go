package repository

import (
	"testing"

	"github.com/JoaoVitor615/bedroom-telemetry-system/internal/domain"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestMongoAdapter_ImplementsIDoorRepository(t *testing.T) {
	// Compile-time assertion that *MongoAdapter implements domain.IDoorRepository
	var _ domain.IDoorRepository = (*MongoAdapter)(nil)
}

func TestNewMongoAdapter(t *testing.T) {
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("failed to create mongo client: %v", err)
	}

	db := client.Database("test_db")
	adapter := NewMongoAdapter(db, "door_events")

	if adapter == nil {
		t.Fatal("expected NewMongoAdapter to return a non-nil repository instance")
	}

	mongoAdapter, ok := adapter.(*MongoAdapter)
	if !ok {
		t.Fatal("expected adapter to be of type *MongoAdapter")
	}

	if mongoAdapter.collection == nil {
		t.Fatal("expected mongoAdapter collection to be initialized")
	}
}
