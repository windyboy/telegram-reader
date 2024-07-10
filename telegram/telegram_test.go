package telegram

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gzzn.com/airport/serial/config"
	"gzzn.com/airport/serial/logger"
)

func TestGetTelegramSequence(t *testing.T) {
	Convey("Given a telegram string", t, func() {
		if err := config.InitParameter(); err != nil {
			t.Fatalf("Failed to initialize config parameter: %v", err)
		}
		parameter := config.GetParameter() // Assuming there's a method to retrieve the loaded parameter
		loggerConfig := logger.LoggerConfig{
			Level:      "debug", // Set the log level to debug or any desired level
			Filename:   "",
			MaxSize:    0,
			MaxBackups: 0,
			MaxAge:     0,
			Compress:   false,
		}
		logger.InitLogger(loggerConfig) // Set logger config

		telegram := `ZCZC TMQ2627 151600


FF ZBTJZXZX


151600 ZGSDZTZX


(DEP-OKA2832/A2426-ZGSD1600-ZBTJ)







NNNN`
		Convey("When the telegram is matched by the sequence pattern", func() {
			seqPattern := parameter.Telegram.SeqTag
			Convey("Then the sequence is returned", func() {
				seq := GetTelegramSequence(telegram, seqPattern)
				So(seq, ShouldEqual, "TMQ2627")
			})
		})
	})
}
