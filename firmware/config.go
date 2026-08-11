package main

import "machine"

const (
	PinSDA = machine.GPIO21
	PinSCL = machine.GPIO22

	// Digital Pin for MC-38 Magnetic Door Sensor
	PinDoor = machine.GPIO4

	// Polling and telemetry intervals (in milliseconds)
	TempTelemetryIntervalMs = 2000
	DoorPollingIntervalMs   = 50
)
