package telegram

import (
	"bytes"
	"regexp"
	"sync"

	"gzzn.com/airport/serial/config"
)

const SequenceUnknow = "TMQ----"

var (
	buffer          bytes.Buffer // Buffer to store incoming data
	mu              sync.Mutex   // Mutex to protect buffer access
	telegramConfig  config.TelegramConfig
	patternEndTag   *regexp.Regexp
	patternSeqTag   *regexp.Regexp
	patternTelegram *regexp.Regexp
)

// Init initializes the telegram package.
// It retrieves the telegram configuration from the parameter and sets up the end tag pattern.
// It also initializes the logger.
func Init(parameter *config.Parameter) {
	if patternEndTag != nil {
		return
	}
	// Retrieve the telegram configuration from the parameter
	telegramConfig = parameter.Telegram

	// Set up the end tag pattern using the configuration
	patternEndTag = regexp.MustCompile(telegramConfig.EndTag)

	// Set up the sequence tag pattern using the configuration
	patternSeqTag = regexp.MustCompile(telegramConfig.SeqTag)

	// Set up the telegram pattern using the configuration
	patternTelegram = regexp.MustCompile(telegramConfig.PatternSplit)

}

func Append(b byte) []string {
	mu.Lock()
	defer mu.Unlock()

	buffer.WriteByte(b)
	// currentBuffer := buffer.String()

	// sugar.Debugf("Buffer: %s", currentBuffer)

	if telegrams := getTelegrams(buffer.String()); len(telegrams) > 0 {
		buffer.Reset() // Reset the buffer if a match is found
		// logger.GetLogger().Debugf("Got %d telegrams", len(telegrams))
		// fmt.Println("len :", len(telegrams))
		return telegrams
	}

	return nil
}

// GetTelegramSequence extracts and returns the telegram sequence from the given telegram
// using the provided sequence pattern. If no sequence is found, it returns an empty string.
func GetSequence(telegram string) string {

	if match := patternSeqTag.FindStringSubmatch(telegram); len(match) > 1 {
		// logger.GetLogger().Debugf("Matched telegram sequence: %s", match[1])
		return match[1]
	}

	return SequenceUnknow
}

// processData checks if the given data matches the provided telegram end tag.
// If a match is found, it returns the data, otherwise returns an empty string.
func getTelegrams(data string) []string {
	if patternEndTag.MatchString(data) {
		// Split the data based on the end tag pattern
		telegrams := patternTelegram.FindAllString(data, -1)
		// fmt.Println("length :", len(telegrams))
		// fmt.Println(telegrams)
		return telegrams
	}
	return nil
}
