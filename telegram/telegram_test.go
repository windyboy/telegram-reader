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
	Context("Given a telegram string", func() {
		var sugar *zap.SugaredLogger

		BeforeEach(func() {
			logger.Init()
			err := config.InitParameter()
			Expect(err).NotTo(HaveOccurred(), "Failed to initialize config parameter")
			sugar = logger.SugaredLogger()
			SetSugaredLogger(sugar)
		})

		It("should return the sequence when the telegram is matched by the sequence pattern", func() {
			text := `ZCZC TMQ2627 151600


FF ZBTJZXZX


151600 ZGSDZTZX


(DEP-OKA2832/A2426-ZGSD1600-ZBTJ)







NNNN`
			seqPattern := "ZCZC\\s(\\S+)\\s"
			seq := GetTelegramSequence(text, seqPattern)
			Expect(seq).To(Equal("TMQ2627"))
		})
	})
})
