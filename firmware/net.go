package main

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqttClient mqtt.Client

type MqttPublisher struct {
	topic   string
	payload string
	retain  bool
}

func NewMqttPublisher(topic string, retain bool) *MqttPublisher {
	return &MqttPublisher{
		topic:  topic,
		retain: retain,
	}
}

func (p *MqttPublisher) PublishMessage(payload string) {
	if mqttClient == nil || !mqttClient.IsConnected() {
		fmt.Println("⚠️ [MQTT] Cannot publish: client is disconnected")
		return
	}

	token := mqttClient.Publish(p.topic, 0, p.retain, payload)

	token.Wait()

	if token.Error() != nil {
		fmt.Printf("❌ [MQTT] Error publishing to %s: %v\n", p.topic, token.Error())
	} else {
		fmt.Printf("📡 [MQTT PUBLISH] Topic: %s | Payload: %s (Retain: %t)\n", p.topic, payload, p.retain)
	}

}
