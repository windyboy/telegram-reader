package telegram

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gzzn.com/airport/serial/config"
)

// TestTelegram runs the test suite for the Telegram package.
func TestTelegram(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Telegram Suite")
}

var _ = Describe("Telegram Processing", func() {
	// var telegramConfig config.TelegramConfig
	// Setup common to all tests in this suite.
	BeforeEach(func() {
		telegramConfig = config.GetParameter().Telegram
		Init()
	})

	Context("When processing a complete telegram string", func() {
		// Example telegram string.
		const telegramText = `ZCZC TMQ2627 151600
FF ZBTJZXZX
151600 ZGSDZTZX
(DEP-OKA2832/A2426-ZGSD1600-ZBTJ)
NNNN`

		It("Extracts the correct sequence from the telegram", func() {
			// Define the pattern to extract the sequence.
			// Extract the sequence.
			sequence := GetSequence(telegramText)
			// Verify the extracted sequence is as expected.
			Expect(sequence).To(Equal("TMQ2627"), "The extracted sequence should match the expected value.")
		})
	})

	Context("When processing a string containing multiple telegrams", func() {
		// Example string containing two telegrams.
		const telegramsText = `ZCZC TMQ2611 151524
FF ZBTJZPZX
151524 ZGGGZPZX
(ARR-CSN3136/A0006-ZBTJ-ZGGG1521)
NNNN
ZCZC TMQ2609 151524
FF ZBTJZPZX
151523 ZBACZQZX
(ARR-CXA8324/A4001-ZLXY-ZBTJ1522)
NNNN`

		// Expected individual telegrams after splitting.
		const expectedFirstTelegram = `ZCZC TMQ2611 151524
FF ZBTJZPZX
151524 ZGGGZPZX
(ARR-CSN3136/A0006-ZBTJ-ZGGG1521)
NNNN`
		const expectedSecondTelegram = `ZCZC TMQ2609 151524
FF ZBTJZPZX
151523 ZBACZQZX
(ARR-CXA8324/A4001-ZLXY-ZBTJ1522)
NNNN`

		It("Correctly splits the string into individual telegrams", func() {
			// Define the pattern to split the telegrams.
			// var splitPattern = "(?s)ZCZC.*?NNNN"
			// Split the string into telegrams.
			telegrams := getTelegrams(telegramsText)
			// Verify the number of telegrams and their content.
			Expect(telegrams).To(HaveLen(2), "There should be exactly two telegrams.")
			Expect(telegrams[0]).To(Equal(expectedFirstTelegram), "The first telegram should match the expected content.")
			Expect(telegrams[1]).To(Equal(expectedSecondTelegram), "The second telegram should match the expected content.")
		})
	})
})
