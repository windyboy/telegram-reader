package telegram

import (
	"bytes"
	"regexp"
	"sync"

	"go.uber.org/zap"
	"gzzn.com/airport/serial/config"
	"gzzn.com/airport/serial/logger"
)

const SequenceUnknow = "TMQ----"

var (
	buffer          bytes.Buffer // Buffer to store incoming data
	mu              sync.Mutex   // Mutex to protect buffer access
	sugar           *zap.SugaredLogger
	telegramConfig  config.TelegramConfig
	patternEndTag   *regexp.Regexp
	patternSeqTag   *regexp.Regexp
	patternTelegram *regexp.Regexp
	initialized     = false
)

// Init initializes the telegram package.
// It retrieves the telegram configuration from the parameter and sets up the end tag pattern.
// It also initializes the logger.
func Init() {
	if initialized {
		return
	}
	// Retrieve the telegram configuration from the parameter
	telegramConfig = config.GetParameter().Telegram

	// Set up the end tag pattern using the configuration
	patternEndTag = regexp.MustCompile(telegramConfig.EndTag)

	// Set up the sequence tag pattern using the configuration
	patternSeqTag = regexp.MustCompile(telegramConfig.SeqTag)

	// Set up the telegram pattern using the configuration
	patternTelegram = regexp.MustCompile(telegramConfig.PatternSplit)

	// Initialize the logger
	logger.Init()
	sugar = logger.SugaredLogger()

	initialized = true
}

// Append appends the given data to the buffer and processes it to extract telegrams.
// It locks the buffer to ensure thread safety and releases the lock when done.
// If a complete telegram is found in the buffer, it resets the buffer and returns the extracted telegrams.
// Otherwise, it returns nil.
func Append(data string) []string {
	mu.Lock()
	defer mu.Unlock()

	buffer.WriteString(data)
	currentBuffer := buffer.String()

	sugar.Debugf("Buffer: %s", currentBuffer)

	if telegrams := processData(currentBuffer); len(telegrams) > 0 {
		buffer.Reset() // Reset the buffer if a match is found
		sugar.Debugf("Got %d telegrams", len(telegrams))
		return telegrams
	}

	return nil
}

// GetTelegramSequence extracts and returns the telegram sequence from the given telegram
// using the provided sequence pattern. If no sequence is found, it returns an empty string.
func GetTelegramSequence(telegram string) string {

	if match := patternSeqTag.FindStringSubmatch(telegram); len(match) > 1 {
		sugar.Debugf("Matched telegram sequence: %s", match[1])
		return match[1]
	}

	return SequenceUnknow
}

// GetTelegramFromText extracts and returns all telegrams from the given text
// using the provided telegram pattern.
func GetTelegramFromText(text string) []string {
	return patternTelegram.FindAllString(text, -1)
}

// processData checks if the given data matches the provided telegram end tag.
// If a match is found, it returns the data, otherwise returns an empty string.
func processData(data string) []string {
	if patternEndTag.MatchString(data) {
		// Split the data based on the end tag pattern
		telegrams := GetTelegramFromText(data)
		return telegrams
	}
	return nil
}
