package main

import "machine"

const (
	// Wi-Fi Credentials
	WifiSSID = "AP333"
	WifiPass = "ap333ap333"

	// MQTT Broker Credentials
	MqttBroker   = "192.168.1.75"
	MqttPort     = 1883
	MqttClientID = "bedroom-esp32-node"

	// MQTT Topics
	TopicTemperature = "bedroom/temperature"
	TopicDoorEvent   = "bedroom/door/event"

	// Hardware Pins
	PinSDA  = machine.GPIO21
	PinSCL  = machine.GPIO22
	PinDoor = machine.GPIO4

	// Polling and telemetry intervals (in milliseconds)
	TempTelemetryIntervalMs = 2000
	DoorPollingIntervalMs   = 50
)
