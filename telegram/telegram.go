package telegram

import (
	"bytes"
	"regexp"
	"sync"

	"gzzn.com/airport/serial/logger"
)

var (
	buffer bytes.Buffer
	mu     sync.Mutex // Mutex to protect buffer access
)

func Append(data string, telegramEndTag string) string {
	sugar := logger.SugaredLogger()
	mu.Lock() // Lock the mutex
	buffer.WriteString(data)
	currentBuffer := buffer.String()
	mu.Unlock() // Unlock the mutex
	sugar.Debugf("Buffer: %s", currentBuffer)
	re, err := regexp.Compile(telegramEndTag)
	if err != nil {
		sugar.Fatalf("Error compiling telegram pattern: %v", err)
	}
	if match := processData(re, currentBuffer); match != "" {
		sugar.Debugf("Matched telegram: %s", match)
		return match
	}
	return ""
}

func GetTelegramSequence(telegram string, seqPattern string) string {
	sugar := logger.SugaredLogger()
	// sugar.Infof("Sequence pattern: %s", seqPattern)
	re, err := regexp.Compile(seqPattern)
	if err != nil {
		sugar.Fatalf("Error compiling telegram pattern: %v", err)
	}
	match := re.FindStringSubmatch(telegram)
	if len(match) > 1 {
		sugar.Debugf("Matched telegram sequence: %s", match[1])
		return match[1]
	}
	return ""
}

// Function to process the buffer
func processData(re *regexp.Regexp, data string) string {

	if re.MatchString(data) {
		buffer.Reset()
		return data
	}
	return ""
}
