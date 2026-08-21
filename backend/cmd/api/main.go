package main

import (
	"log"
	"net/http"
	"os"

	repoModels "github.com/JoaoVitor615/bedroom-telemetry-system/internal/domain/repository"

	"github.com/JoaoVitor615/bedroom-telemetry-system/internal/adapters/repository"
	"github.com/JoaoVitor615/bedroom-telemetry-system/internal/domain/controller"
	"github.com/JoaoVitor615/bedroom-telemetry-system/internal/domain/service"
	"github.com/JoaoVitor615/bedroom-telemetry-system/internal/mongodb"
	"github.com/JoaoVitor615/bedroom-telemetry-system/internal/mqtt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("🚀 Starting Bedroom API...")

	// 1. Load the variables in .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Warning: .env file not found.")
	}

	// 2. Get credentials
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB_NAME")
	mqttBroker := os.Getenv("MQTT_BROKER_URI")
	apiPort := os.Getenv("PORT")

	// Port fallback
	if apiPort == "" {
		apiPort = "8080"
	}

	// 3. Initialize Domain State
	appState := &repoModels.CurrentState{}

	// 4. Setup Infrastructure (MongoDB Connection)
	mongoClient, err := mongodb.NewClient(mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	db := mongoClient.Database(dbName)

	// 5. Initialize Adapters
	mongoAdapter := repository.NewMongoAdapter(db, "door_events")

	// 6. Initialize Core Domain Services
	bedroomService := service.NewBedroomService(mongoAdapter, appState)

	// 7. Initialize Controllers and Consumers
	httpController := controller.NewHTTPController(bedroomService)

	mqttConsumer := mqtt.NewConsumer(mqttBroker, bedroomService)
	mqttConsumer.StartListening()

	// 8. Setup HTTP Router (Chi)
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Get("/state", httpController.GetCurrentState)
		r.Get("/history", httpController.GetHistory)
	})

	// 9. Start Server
	serverAddress := ":" + apiPort
	log.Printf("🌐 HTTP Server running on port %s\n", serverAddress)
	if err := http.ListenAndServe(serverAddress, r); err != nil {
		log.Fatal(err)
	}
}
