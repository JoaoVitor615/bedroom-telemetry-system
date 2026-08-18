package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/eclipse/paho.golang/paho"
	"tinygo.org/x/drivers/netdev"
	"tinygo.org/x/drivers/netlink"
	link "tinygo.org/x/espradio/netlink"
)

// global handle for the new Paho MQTT Client (usa ponteiro para paho.Client)
var mqttClient *paho.Client

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

// publishMessage sends a payload to a target topic using the new Paho lib
func (p *MqttPublisher) PublishMessage(payload string) {
	// A nova biblioteca gerencia a conexão na camada de rede (net.Conn).
	// Se a rede cair, o método Publish retornará um erro automaticamente.
	if mqttClient == nil {
		fmt.Println("⚠️ [MQTT] Cannot publish: client is nil")
		return
	}

	msg := &paho.Publish{
		Topic:   p.topic,
		Payload: []byte(payload),
		QoS:     0,
		Retain:  p.retain,
	}

	// Envia a mensagem passando um context nativo do Go em vez do antigo "token.Wait()"
	_, err := mqttClient.Publish(context.Background(), msg)

	if err != nil {
		fmt.Printf("❌ [MQTT] Error publishing to %s: %v\n", p.topic, err)
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

	// format target broker address (e.g., 192.168.1.100:1883)
	// Nota: a nova lib usa o dialer padrão TCP, então não precisamos do prefixo "tcp://"
	brokerAddress := fmt.Sprintf("%s:%d", MqttBroker, MqttPort)
	fmt.Printf("🔌 Connecting TCP to %s...\n", brokerAddress)

	// Abre a conexão TCP padrão usando a internet configurada pelo netdev
	conn, err := net.Dial("tcp", brokerAddress)
	if err != nil {
		return fmt.Errorf("failed to open TCP connection: %w", err)
	}

	// instantiate new Paho client
	mqttClient = paho.NewClient(paho.ClientConfig{
		Conn: conn,
		// Substitui o OnConnectionLost antigo
		OnClientError: func(err error) {
			fmt.Printf("⚠️ Connection lost to MQTT Broker: %v\n", err)
		},
	})

	// Setup do pacote de conexão MQTT
	connectPacket := &paho.Connect{
		ClientID:   MqttClientID,
		CleanStart: true,
		KeepAlive:  30,
	}

	fmt.Println("🔄 Authenticating MQTT...")

	// Cria um contexto com timeout para a tentativa de conexão (substitui o token.Wait)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Envia o pacote de Connect e aguarda o Acknowledgement do Broker
	connAck, err := mqttClient.Connect(ctx, connectPacket)
	if err != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %w", err)
	}

	if connAck.ReasonCode != 0 {
		return fmt.Errorf("broker rejected connection with reason code: %d", connAck.ReasonCode)
	}

	// Substitui o OnConnect antigo
	fmt.Println("✅ Successfully connected to MQTT Broker!")

	return nil
}
