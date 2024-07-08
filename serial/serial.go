package serial

import (
	"fmt"

	serial "go.bug.st/serial"
	"gzzn.com/airport/serial/logger"
)

func ReadFromPort(mode *serial.Mode, portName string, bufferSize int, dataChannel chan<- []byte) error {
	sugar := logger.SugaredLogger()

	sugar.Infof("Opening port: %s with mode: %+v", portName, mode)

	port, err := serial.Open(portName, mode)
	if err != nil {
		return fmt.Errorf("error opening port: %w", err)
	}
	defer port.Close()

	buf := make([]byte, bufferSize)
	for {
		num, err := port.Read(buf)
		if err != nil {
			return fmt.Errorf("error reading from port: %w", err)
		}
		readData := buf[:num]
		// sugar.Debugf("Read %d bytes: %s", num, string(readData))

		// Log the data being sent over the channel
		// sugar.Debug("Sending data over channel:", string(readData))
		dataChannel <- readData
	}
}
