package telegram

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gzzn.com/airport/serial/config"
	"gzzn.com/airport/serial/logger"
)

func TestGetTelegramSequence(t *testing.T) {
	Convey("Given a telegram string", t, func() {

		parameter, _ := config.LoadConfigFromEnv()
		logger.SetParameter(&parameter)
		logger.InitTestLogger()

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
