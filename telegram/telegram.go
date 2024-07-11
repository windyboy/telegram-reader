package telegram

import (
	"bytes"
	"regexp"
	"sync"

	"go.uber.org/zap"
	"gzzn.com/airport/serial/config"
	"gzzn.com/airport/serial/logger"
)

var (
	buffer         bytes.Buffer // Buffer to store incoming data
	mu             sync.Mutex   // Mutex to protect buffer access
	sugar          *zap.SugaredLogger
	telegramConfig config.TelegramConfig
	patternEndTag  *regexp.Regexp
)

func Init() {
	if param, err := config.GetParameter(); err != nil {
		panic("Error getting parameter")
	} else {
		telegramConfig = param.Telegram
		patternEndTag = regexp.MustCompile(telegramConfig.EndTag)
	}
	logger.Init()
	sugar = logger.SugaredLogger()
}

// Append adds data to the buffer and checks for matches against the telegram end tag.
// If a match is found, it returns the matched telegram, otherwise returns an empty string.
func Append(data string) string {

	mu.Lock()
	defer mu.Unlock()

	buffer.WriteString(data)
	currentBuffer := buffer.String()

	sugar.Debugf("Buffer: %s", currentBuffer)

	if match := processData(currentBuffer); match != "" {
		buffer.Reset() // Reset the buffer if a match is found
		sugar.Debugf("Matched telegram: %s", match)
		return match
	}

	return ""
}

// GetTelegramSequence extracts and returns the telegram sequence from the given telegram
// using the provided sequence pattern. If no sequence is found, it returns an empty string.
func GetTelegramSequence(telegram string, seqPattern string) string {

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

// GetTelegramFromText extracts and returns all telegrams from the given text
// using the provided telegram pattern.
func GetTelegramFromText(text string, pattern string) []string {

	re, err := regexp.Compile(pattern)
	if err != nil {
		sugar.Fatalf("Error compiling telegram pattern: %v", err)
	}

	matches := re.FindAllString(text, -1)
	return matches
}

// processData checks if the given data matches the provided telegram end tag.
// If a match is found, it returns the data, otherwise returns an empty string.
func processData(data string) string {
	if patternEndTag.MatchString(data) {
		return data
	}
	return ""
}
