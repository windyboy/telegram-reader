# Serial Port Civial Authority Aviation Telegram Message Receiver

## Description

This is a simple application that reads telegram messages from a serial port and publishes them to a NATS server.
This project is a specialized serial port reader designed for civil aviation applications. It captures telegram messages from serial ports and forwards them to a NATS messaging system, making the data available for broader aviation systems integration.

Key Features:
- Serial port data acquisition with configurable parameters
- Aviation telegram message parsing and validation
- Real-time message publishing via NATS
- Prometheus metrics for monitoring
- Comprehensive logging system
- Flexible configuration through TOML, CLI flags, or environment variables

The application serves as a critical bridge between serial-based aviation equipment and modern distributed systems, ensuring reliable data transmission while maintaining operational monitoring capabilities.

## Usage