package main

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqttClient mqtt.Client

func PublishMessage(topic string, payload string, retain bool) {
	if mqttClient == nil || !mqttClient.IsConnected() {
		fmt.Println("⚠️ [MQTT] Cannot publish: client is disconnected")
		return
	}

	token := mqttClient.Publish(topic, 0, retain, payload)

	token.Wait()

	if token.Error() != nil {
		fmt.Printf("❌ [MQTT] Error publishing to %s: %v\n", topic, token.Error())
	} else {
		fmt.Printf("📡 [MQTT PUBLISH] Topic: %s | Payload: %s (Retain: %t)\n", topic, payload, retain)
	}

}
