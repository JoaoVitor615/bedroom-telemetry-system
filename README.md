# 🏠 Room Telemetry System

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Firmware](https://img.shields.io/badge/firmware-ESP32%20%7C%20PlatformIO-orange)
![Backend](https://img.shields.io/badge/backend-Go%20%7C%20Hexagonal-00ADD8)
![Mobile](https://img.shields.io/badge/mobile-Swift%20%7C%20iOS-F05138)

A lightweight, end-to-end IoT system for real-time bedroom environment monitoring. It reads temperature, humidity, and door status via an ESP32, streams telemetry over MQTT, persists door event logs in MongoDB through a Go backend, and presents live data in a native iOS app.

## 📐 System Architecture

The project is structured as a **Monorepo** and follows a highly decoupled design, using Hexagonal Architecture (Ports & Adapters) for the backend.

```text
[ ESP32 + Sensors ] ──(MQTT)──> [ Go API ] ──(REST HTTP)──> [ iOS App (Swift) ]
                                    │
                               [ MongoDB ]
