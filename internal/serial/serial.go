package serial

import (
	"fmt"

	serial "go.bug.st/serial"
	"go.uber.org/zap"
)

func ReadFromPort(mode *serial.Mode, portName string) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	sugar.Infof("Opening port: %s with mode: %+v", portName, mode)

	port, err := serial.Open(portName, mode)
	if err != nil {
		return fmt.Errorf("error opening port: %w", err)
	}
	defer port.Close()

	buf := make([]byte, 100)
	num, err := port.Read(buf)
	if err != nil {
		return fmt.Errorf("error reading from port: %w", err)
	}

	sugar.Infof("Read %d bytes: %s", num, string(buf[:num]))
	return nil
}
