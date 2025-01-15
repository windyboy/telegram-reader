package telegram

import (
	"bytes"
	"regexp"
	"sync"
)

const SequenceUnknow = "TMQ----"

var (
	buffer          bytes.Buffer // Buffer to store incoming data
	mu              sync.Mutex   // Mutex to protect buffer access
	patternEndTag   = regexp.MustCompile("NNNN")
	patternSeqTag   = regexp.MustCompile(`ZCZC\s(\S+)\s`)
	patternTelegram = regexp.MustCompile(`(?s)ZCZC.*?NNNN`)
)

// Init initializes the telegram package.
// func Init(_ interface{}) {
// 	// No initialization needed since patterns are hardcoded
// 	return
// }

func Append(b byte) []string {
	mu.Lock()
	defer mu.Unlock()

	buffer.WriteByte(b)

	if telegrams := getTelegrams(buffer.String()); len(telegrams) > 0 {
		buffer.Reset() // Reset the buffer if a match is found
		return telegrams
	}

	return nil
}

// GetSequence extracts and returns the telegram sequence from the given telegram.
// If no sequence is found, it returns SequenceUnknow.
func GetSequence(telegram string) string {
	if match := patternSeqTag.FindStringSubmatch(telegram); len(match) > 1 {
		return match[1]
	}

	return SequenceUnknow
}

// getTelegrams checks if the given data matches the telegram end tag.
// If matches are found, returns the complete telegrams, otherwise returns nil.
func getTelegrams(data string) []string {
	if patternEndTag.MatchString(data) {
		return patternTelegram.FindAllString(data, -1)
	}
	return nil
}
