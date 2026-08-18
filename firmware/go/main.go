package main

import (
	"fmt"
	"machine"
	"time"

	"tinygo.org/x/drivers/bme280"
)

func main() {
	// wait 2 seconds for hardware stabilization
	time.Sleep(2 * time.Second)
	fmt.Println("🚀 Starting bedroom monitoring system...")

	err := InitNetwork()
	if err != nil {
		fmt.Printf("❌ Network initialization failed: %v\n", err)
		fmt.Println("⚠️ Running in offline mode (local serial logs only)...")
	}

	// retain is FALSE for door events (historical state log)
	doorPublisher := NewMqttPublisher(TopicDoorEvent, false)
	// retain is TRUE for temperature (keeps last known state for immediate UI consumption)
	tempPublisher := NewMqttPublisher(TopicTemperature, true)

	// configure I2C for BME280
	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: machine.KHz * 400,
		SCL:       PinSCL,
		SDA:       PinSDA,
	})

	sensorBME := bme280.New(machine.I2C0)
	sensorBME.ConfigureWithSettings(bme280.Config{Mode: bme280.ModeNormal})

	PinDoor.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	// background Goroutine: Door Monitor (MC-38 Sensor)
	go func() {
		lastDoorState := PinDoor.Get()

		for {
			currentDoorState := PinDoor.Get()

			// Trigger event only when physical state changes
			if currentDoorState != lastDoorState {
				lastDoorState = currentDoorState

				if currentDoorState {
					fmt.Println("🚪 [EVENT] Door was OPENED")
					doorPublisher.PublishMessage(`{"event":"OPENED"}`)
				} else {
					fmt.Println("🚪 [EVENT] Door was CLOSED")
					doorPublisher.PublishMessage(`{"event":"CLOSED"}`)
				}
			}

			time.Sleep(DoorPollingIntervalMs * time.Millisecond)
		}
	}()

	// temperature & humidity Monitor (BME280)
	for {
		rawTemp, err := sensorBME.ReadTemperature()
		if err == nil {
			tempCelsius := float64(rawTemp) / 1000.0
			rawHumidity, _ := sensorBME.ReadHumidity()
			humidity := float64(rawHumidity) / 100.0

			fmt.Printf("🌡️ [TELEMETRY] Temp: %.2f °C | Humidity: %.2f %%\n", tempCelsius, humidity)
			payload := fmt.Sprintf(`{"temperature":%.2f,"humidity":%.2f}`, tempCelsius, humidity)
			tempPublisher.PublishMessage(payload)
		} else {
			fmt.Println("⚠️ [ERROR] Failed to read from BME280 sensor:", err)
		}

		time.Sleep(TempTelemetryIntervalMs * time.Millisecond)
	}

}
