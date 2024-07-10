package telegram

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"gzzn.com/airport/serial/config"
	"gzzn.com/airport/serial/logger"
)

func TestTelegram(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Telegram Suite")
}

var _ = Describe("Telegram", func() {
	Context("Telegram test mode", func() {
		var sugar *zap.SugaredLogger

		BeforeEach(func() {
			logger.Init()
			err := config.InitParameter()
			Expect(err).NotTo(HaveOccurred(), "Failed to initialize config parameter")
			sugar = logger.SugaredLogger()
			SetSugaredLogger(sugar)
		})

		Context("A complete telegram string", func() {
			text := `ZCZC TMQ2627 151600


FF ZBTJZXZX


151600 ZGSDZTZX


(DEP-OKA2832/A2426-ZGSD1600-ZBTJ)







NNNN`

			It("should return the sequence when the telegram is matched by the sequence pattern", func() {

				seqPattern := "ZCZC\\s(\\S+)\\s"
				seq := GetTelegramSequence(text, seqPattern)
				Expect(seq).To(Equal("TMQ2627"))
			})
		}) //End of Context

		Context("A string with 2 telegrams", func() {
			text := `
ZCZC TMQ2611 151524


FF ZBTJZPZX


151524 ZGGGZPZX


(ARR-CSN3136/A0006-ZBTJ-ZGGG1521)







NNNN

ZCZC TMQ2609 151524


FF ZBTJZPZX


151523 ZBACZQZX


(ARR-CXA8324/A4001-ZLXY-ZBTJ1522)







NNNN
			`
			pattern := "(?s)ZCZC.*?NNNN"
			telegrams := GetTelegramFromText(text, pattern)
			t1 := `ZCZC TMQ2611 151524


FF ZBTJZPZX


151524 ZGGGZPZX


(ARR-CSN3136/A0006-ZBTJ-ZGGG1521)







NNNN`
			t2 := `ZCZC TMQ2609 151524


FF ZBTJZPZX


151523 ZBACZQZX


(ARR-CXA8324/A4001-ZLXY-ZBTJ1522)







NNNN`
			It("should return 2 telegrams", func() {
				Expect(len(telegrams)).To(Equal(2))
				Expect(telegrams[0]).To(Equal(t1))
				Expect(telegrams[1]).To(Equal(t2))

			})

		}) //End of Context

	}) //End of Describe
}) //End of Context
