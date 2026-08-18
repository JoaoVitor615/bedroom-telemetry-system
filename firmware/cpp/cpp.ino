#include <WiFi.h>
#include <PubSubClient.h>
#include <Wire.h>
#include <Adafruit_Sensor.h>
#include <Adafruit_BME280.h>
#include "config.h"

// Global Instances
WiFiClient espClient;
PubSubClient mqtt(espClient);
Adafruit_BME280 bme;

// State Variables
unsigned long lastTempRead = 0;
int lastDoorState = -1;

void setupWiFi() {
    Serial.printf("📶 Connecting to Wi-Fi SSID: %s...\n", WIFI_SSID);
    WiFi.begin(WIFI_SSID, WIFI_PASS);

    while (WiFi.status() != WL_CONNECTED) {
        delay(500);
        Serial.print(".");
    }
    Serial.printf("\n✅ Wi-Fi Connected! IP: %s\n", WiFi.localIP().toString().c_str());
}

void reconnectMQTT() {
    while (!mqtt.connected()) {
        Serial.printf("🔌 Connecting to MQTT Broker at %s:%d...\n", MQTT_BROKER, MQTT_PORT);
        
        if (mqtt.connect(MQTT_CLIENT_ID)) {
            Serial.println("✅ Successfully connected to MQTT Broker!");
        } else {
            Serial.printf("⚠️ Connection failed, state: %d. Retrying in 5 seconds...\n", mqtt.state());
            delay(5000);
        }
    }
}

void publishMQTT(const char* topic, const char* payload) {
    if (mqtt.connected()) {
        mqtt.publish(topic, payload, false); // false = no retain
        Serial.printf("📡 [MQTT PUBLISH] Topic: %s | Payload: %s\n", topic, payload);
    }
}

void setup() {
    Serial.begin(115200);
    delay(2000); // wait for hardware stabilization
    Serial.println("🚀 Starting bedroom monitoring system...");

    // Door Pin Setup
    pinMode(PIN_DOOR, INPUT_PULLUP);
    lastDoorState = digitalRead(PIN_DOOR);

    // BME280 Sensor Setup via I2C
    Wire.begin(PIN_SDA, PIN_SCL);
    // Note: The default BME280 address is usually 0x76. If it fails, try 0x77.
    if (!bme.begin(0x76, &Wire)) { 
        Serial.println("⚠️ [ERROR] Failed to read from BME280 sensor!");
    }

    setupWiFi();
    mqtt.setServer(MQTT_BROKER, MQTT_PORT);
}

void loop() {
    // 1. Keep the MQTT connection alive
    if (!mqtt.connected()) {
        reconnectMQTT();
    }
    mqtt.loop(); // Essential: processes incoming messages and keep-alive

    unsigned long now = millis();

    // 2. Door Monitor (similar to the Goroutine)
    int currentDoorState = digitalRead(PIN_DOOR);
    
    if (currentDoorState != lastDoorState) {
        lastDoorState = currentDoorState;

        // If using PULLUP with the switch to GND: HIGH = Open, LOW = Closed.
        if (currentDoorState == HIGH) {
            Serial.println("🚪 [EVENT] Door was OPENED");
            publishMQTT(TOPIC_DOOR_EVENT, "OPENED");
        } else {
            Serial.println("🚪 [EVENT] Door was CLOSED");
            publishMQTT(TOPIC_DOOR_EVENT, "CLOSED");
        }
        delay(DOOR_POLLING_INTERVAL_MS); // Simple debounce
    }

    // 3. Temperature & Humidity Monitor
    if (now - lastTempRead >= TEMP_TELEMETRY_INTERVAL_MS) {
        lastTempRead = now;
        
        float tempCelsius = bme.readTemperature();
        float humidity = bme.readHumidity();

        // Check if the reading was valid
        if (!isnan(tempCelsius) && !isnan(humidity)) {
            Serial.printf("🌡️ [TELEMETRY] Temp: %.2f °C | Humidity: %.2f %%\n", tempCelsius, humidity);
            
            // Convert floats to string to send via MQTT
            char tempStr[10];
            char humStr[10];
            dtostrf(tempCelsius, 4, 2, tempStr);
            dtostrf(humidity, 4, 2, humStr);

            publishMQTT(TOPIC_TEMPERATURE, tempStr);
            publishMQTT("bedroom/humidity", humStr);
        } else {
            Serial.println("⚠️ [ERROR] Failed to read from BME280 sensor");
        }
    }
}