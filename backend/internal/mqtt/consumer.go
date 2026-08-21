package mqtt

import (
	"context"
	"fmt"
	"log"

	"github.com/JoaoVitor615/bedroom-telemetry-system/internal/domain"
	paho "github.com/eclipse/paho.mqtt.golang"
)

type Consumer struct {
	client  paho.Client
	service domain.IBedroomService
}

// NewConsumer initializes the MQTT client and injects the domain service
func NewConsumer(brokerURI string, service domain.IBedroomService) *Consumer {
	opts := paho.NewClientOptions().AddBroker(brokerURI).SetClientID("bedroom-api-consumer")

	client := paho.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}

	return &Consumer{
		client:  client,
		service: service,
	}
}

// StartListening subscribes to topics and forwards data to the domain service
func (c *Consumer) StartListening() {
	c.client.Subscribe("bedroom/temperature", 0, func(client paho.Client, msg paho.Message) {
		temp := string(msg.Payload())
		c.service.ProcessTemperature(temp)
		fmt.Printf("🌡️ New Temperature: %s °C\n", temp)
	})

	c.client.Subscribe("bedroom/door/event", 0, func(client paho.Client, msg paho.Message) {
		state := string(msg.Payload())

		go func() {
			err := c.service.ProcessDoorEvent(context.Background(), state)
			if err != nil {
				log.Printf("Failed to process door event: %v\n", err)
			}
		}()
		fmt.Printf("🚪 Door State: %s\n", state)
	})
}
