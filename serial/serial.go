package serial

import (
	"fmt"

	serial "go.bug.st/serial"
	"gzzn.com/airport/serial/logger"
)

// ReadFromPort reads data from a serial port and sends it over a channel.
// It takes the serial mode, port name, buffer size, and a data channel as input.
// It returns an error if there's any issue with opening the port or reading from it.
func ReadFromPort(mode *serial.Mode, portName string, bufferSize int, dataChannel chan<- []byte) error {
	logger := logger.GetLogger()

	logger.Infof("Opening port: %s with mode: %+v", portName, mode)

	port, err := serial.Open(portName, mode)
	if err != nil {
		return fmt.Errorf("error opening port: %w", err)
	}
	defer port.Close()

	buffer := make([]byte, bufferSize)
	for {
		numBytesRead, err := port.Read(buffer)
		if err != nil {
			return fmt.Errorf("error reading from port: %w", err)
		}
		readData := buffer[:numBytesRead]

		// Uncomment the following line if you want to log the read data
		// logger.Debugf("Read %d bytes: %s", numBytesRead, string(readData))

		// Log the data being sent over the channel
		// logger.Debug("Sending data over channel:", string(readData))
		dataChannel <- readData
	}
}
