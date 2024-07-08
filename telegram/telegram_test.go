package telegram

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetTelegramSequence(t *testing.T) {
	Convey("Given a telegram string", t, func() {
		telegram := `ZCZC TMQ2627 151600


FF ZBTJZXZX


151600 ZGSDZTZX


(DEP-OKA2832/A2426-ZGSD1600-ZBTJ)







NNNN
		`
		Convey("When the telegram is matched by the sequence pattern", func() {
			seqPattern := "ZCZC\\s(\\S+)\\s"
			Convey("Then the sequence is returned", func() {
				seq := GetTelegramSequence(telegram, seqPattern)
				So(seq, ShouldEqual, "TMQ2627")
			})
		})
	})
}
