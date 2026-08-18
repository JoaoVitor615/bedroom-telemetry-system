#pragma once

// ==========================================
// NETWORK SETTINGS
// ==========================================
const char* WIFI_SSID = "YOUR_WIFI_SSID";
const char* WIFI_PASS = "YOUR_WIFI_PASSWORD";

// ==========================================
// MQTT SETTINGS
// ==========================================
const char* MQTT_BROKER    = "192.168.0.100"; // Replace with your Broker IP
const int   MQTT_PORT      = 1883;
const char* MQTT_CLIENT_ID = "bedroom-monitor-01";

// ==========================================
// HARDWARE PINS
// ==========================================
const int PIN_DOOR = 4;  // Change to the pin where the door sensor is connected
const int PIN_SDA  = 21; // Default ESP32 I2C
const int PIN_SCL  = 22; // Default ESP32 I2C

// ==========================================
// INTERVALS (in milliseconds)
// ==========================================
const unsigned long TEMP_TELEMETRY_INTERVAL_MS = 5000;
const unsigned long DOOR_POLLING_INTERVAL_MS   = 200;