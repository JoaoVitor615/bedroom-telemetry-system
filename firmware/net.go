package main

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"tinygo.org/x/drivers/netdev"
	"tinygo.org/x/drivers/netlink"
	link "tinygo.org/x/espradio/netlink"
)

// global handle for the Paho MQTT Client
var mqttClient mqtt.Client

type MqttPublisher struct {
	topic   string
	payload string
	retain  bool
}

// newMqttPublisher return a new MQTT publisher
func NewMqttPublisher(topic string, retain bool) *MqttPublisher {
	return &MqttPublisher{
		topic:  topic,
		retain: retain,
	}
}

// publishMessage sends a payload to a target topic using Paho
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

// InitNetwork initializes the native Wi-Fi stack and establishes an MQTT session
func InitNetwork() error {
	fmt.Printf("📶 Connecting to Wi-Fi SSID: %s...\n", WifiSSID)

	// Instantiate the hardware link driver for the ESP32 Wi-Fi radio
	espLink := &link.Esplink{}

	// Register this physical radio interface as TinyGo's default network device (netdev).
	// This maps low-level socket calls (net.Dial) to the ESP32 hardware radio.
	netdev.UseNetdev(espLink)

	// Attempt to authenticate and connect to the local Wi-Fi Access Point
	err := espLink.NetConnect(&netlink.ConnectParams{
		Ssid:       WifiSSID,
		Passphrase: WifiPass,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to Wi-Fi: %w", err)
	}

	fmt.Println("✅ Wi-Fi Connected!")

	// paho MQTT client configuration
	opts := mqtt.NewClientOptions()

	// format target broker address (e.g., tcp://192.168.1.100:1883)
	brokerURL := fmt.Sprintf("tcp://%s:%d", MqttBroker, MqttPort)
	opts.AddBroker(brokerURL)

	opts.SetClientID(MqttClientID)
	opts.SetCleanSession(true)

	// Automatic background reconnection strategy if connection drops
	opts.SetAutoReconnect(true)
	opts.SetConnectRetryInterval(5 * time.Second)

	// Connection event log callbacks
	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("✅ Successfully connected/reconnected to MQTT Broker!")
	}
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		fmt.Printf("⚠️ Connection lost to MQTT Broker: %v\n", err)
	}

	// instantiate and connect
	mqttClient = mqtt.NewClient(opts)

	fmt.Printf("🔌 Connecting to MQTT Broker at %s...\n", brokerURL)
	token := mqttClient.Connect()

	// Block until the initial MQTT TCP handshake completes or fails
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}

	return nil
}
