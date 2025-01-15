package telegram

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Telegram Processing", func() {
	Context("When processing a complete telegram string", func() {
		const telegramText = `ZCZC TMQ2627 151600
FF ZBTJZXZX
151600 ZGSDZTZX
(DEP-OKA2832/A2426-ZGSD1600-ZBTJ)
NNNN`

		It("Extracts the correct sequence from the telegram", func() {
			sequence := GetSequence(telegramText)
			Expect(sequence).To(Equal("TMQ2627"), "The extracted sequence should match the expected value.")
		})
	})

	Context("When processing a string containing multiple telegrams", func() {
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

		It("Correctly splits the string into individual telegrams", func() {
			telegrams := getTelegrams(telegramsText)
			Expect(telegrams).To(HaveLen(2), "There should be exactly two telegrams.")
			Expect(telegrams[0]).To(Equal(`ZCZC TMQ2611 151524
FF ZBTJZPZX
151524 ZGGGZPZX
(ARR-CSN3136/A0006-ZBTJ-ZGGG1521)
NNNN`), "The first telegram should match the expected content.")
			Expect(telegrams[1]).To(Equal(`ZCZC TMQ2609 151524
FF ZBTJZPZX
151523 ZBACZQZX
(ARR-CXA8324/A4001-ZLXY-ZBTJ1522)
NNNN`), "The second telegram should match the expected content.")
		})
	})

	Context("When processing an empty telegram string", func() {
		const emptyText = ""

		It("Returns an empty slice", func() {
			telegrams := getTelegrams(emptyText)
			Expect(telegrams).To(BeEmpty(), "The result should be an empty slice.")
		})
	})

	Context("When processing a malformed telegram string", func() {
		const malformedText = `ZCZC TMQ2627 151600
FF ZBTJZXZX
NNNN`

		It("Handles the malformed telegram gracefully", func() {
			telegrams := getTelegrams(malformedText)
			Expect(telegrams).To(HaveLen(1), "There should be one telegram despite being malformed.")
			Expect(telegrams[0]).To(ContainSubstring("TMQ2627"), "The telegram should still contain the sequence.")
		})
	})

	Context("When processing a telegram with missing end tag", func() {
		const missingEndTagText = `ZCZC TMQ2627 151600
FF ZBTJZXZX
151600 ZGSDZTZX
(DEP-OKA2832/A2426-ZGSD1600-ZBTJ)`

		It("Returns an empty slice", func() {
			telegrams := getTelegrams(missingEndTagText)
			Expect(telegrams).To(BeEmpty(), "The result should be an empty slice due to missing end tag.")
		})
	})
})

// Run the tests
func TestTelegram(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Telegram Suite")
}
