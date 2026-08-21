package main

import (
	"log"
	"net/http"

	repoModels "github.com/JoaoVitor615/bedroom-telemetry-system/internal/domain/repository"

	"github.com/JoaoVitor615/bedroom-telemetry-system/internal/adapters/repository"
	"github.com/JoaoVitor615/bedroom-telemetry-system/internal/domain/controller"
	"github.com/JoaoVitor615/bedroom-telemetry-system/internal/domain/service"
	"github.com/JoaoVitor615/bedroom-telemetry-system/internal/mongodb"
	"github.com/JoaoVitor615/bedroom-telemetry-system/internal/mqtt"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	log.Println("🚀 Starting Bedroom API...")

	// 1. Initialize Domain State
	appState := &repoModels.CurrentState{}

	// 2. Setup Infrastructure (MongoDB Connection)
	mongoClient, err := mongodb.NewClient("mongodb://localhost:27017")
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	db := mongoClient.Database("bedroom_db")

	// 3. Initialize Adapters
	mongoAdapter := repository.NewMongoAdapter(db, "door_events")

	// 4. Initialize Core Domain Services (Injecting Adapters and State)
	bedroomService := service.NewBedroomService(mongoAdapter, appState)

	// 5. Initialize Controllers and Consumers (Injecting Services)
	httpController := controller.NewHTTPController(bedroomService)

	mqttConsumer := mqtt.NewConsumer("tcp://192.168.1.75:1883", bedroomService)
	mqttConsumer.StartListening()

	// 6. Setup HTTP Router (Chi)
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Get("/state", httpController.GetCurrentState)
		r.Get("/history", httpController.GetHistory)
	})

	// 7. Start Server
	log.Println("🌐 HTTP Server running on port :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
