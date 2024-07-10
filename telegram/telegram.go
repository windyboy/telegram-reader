package telegram

import (
	"bytes"
	"regexp"
	"sync"

	"go.uber.org/zap"
	"gzzn.com/airport/serial/logger" // TODO: remove
)

var (
	buffer bytes.Buffer // Buffer to store incoming data
	mu     sync.Mutex   // Mutex to protect buffer access
	sugar  *zap.SugaredLogger
)

// SetSugaredLogger allows the user to set a custom sugared logger
func SetSugaredLogger(sugaredLogger *zap.SugaredLogger) {
	sugar = sugaredLogger
}

// ensureLogger checks if the sugared logger is initialized, and initializes it if not.
func ensureLogger() {
	if sugar == nil {
		SetSugaredLogger(logger.SugaredLogger())
	}
}

// Append adds data to the buffer, checks for matches against the telegram end tag,
// and returns the matched telegram if found.
func Append(data string, telegramEndTag string) string {
	ensureLogger()
	mu.Lock()
	buffer.WriteString(data)
	currentBuffer := buffer.String()
	mu.Unlock()

	sugar.Debugf("Buffer: %s", currentBuffer)

	if match := processData(telegramEndTag, currentBuffer); match != "" {
		sugar.Debugf("Matched telegram: %s", match)
		return match
	}

	return ""
}

// GetTelegramSequence extracts and returns the telegram sequence from the given telegram
// using the provided sequence pattern. If no sequence is found, it returns an empty string.
func GetTelegramSequence(telegram string, seqPattern string) string {
	ensureLogger()
	re, err := regexp.Compile(seqPattern)
	if err != nil {
		sugar.Fatalf("Error compiling sequence pattern: %v", err)
	}

	match := re.FindStringSubmatch(telegram)
	if len(match) > 1 {
		sugar.Debugf("Matched telegram sequence: %s", match[1])
		return match[1]
	}

	return ""
}

// processData checks if the given data matches the provided telegram end tag.
// If a match is found, it resets the buffer and returns the data, otherwise returns an empty string.
func processData(telegramEndTag string, data string) string {
	if isTelegramEndTagMatched(telegramEndTag, data) {
		buffer.Reset() // Reset the buffer if a match is found
		return data
	}
	return ""
}

// isTelegramEndTagMatched checks if the given data matches the provided telegram end tag pattern.
func isTelegramEndTagMatched(telegramEndTag string, data string) bool {
	re, err := regexp.Compile(telegramEndTag)
	if err != nil {
		sugar.Fatalf("Error compiling telegram end tag pattern: %v", err)
	}

	return re.MatchString(data)
}
