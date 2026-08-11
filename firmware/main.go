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

	// configure I2C for BME280
	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: machine.KHz * 400,
		SCL:       PinSCL,
		SDA:       PinSDA,
	})

	sensorBME := bme280.New(machine.I2C0)
	sensorBME.ConfigureWithSettings(bme280.Config{Mode: bme280.ModeNormal})

	PinDoor.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	go func() {
		lastDoorState := PinDoor.Get()

		for {
			currentDoorState := PinDoor.Get()

			// Trigger event only when physical state changes
			if currentDoorState != lastDoorState {
				lastDoorState = currentDoorState

				if currentDoorState {
					fmt.Println("🚪 [EVENT] Door was OPENED")
					// TODO: publishMQTT("bedroom/door/event", "OPENED")
				} else {
					fmt.Println("🚪 [EVENT] Door was CLOSED")
					// TODO: publishMQTT("bedroom/door/event", "CLOSED")
				}
			}

			time.Sleep(DoorPollingIntervalMs * time.Millisecond)
		}
	}()

	for {
		rawTemp := 
	}

}
