package har

import (
	"encoding/json"
	"fmt"

	"github.com/fedesog/webdriver"
	"github.com/jordanpotter/har"
	"github.com/pkg/errors"
)

func New(logEntries []webdriver.LogEntry) (*har.HAR, error) {
	for _, entry := range logEntries {
		var chromeEntry ChromeLogEntry
		if err := json.Unmarshal([]byte(entry.Message), &chromeEntry); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal log entry")
		}
		fmt.Println("TODO", chromeEntry)
	}

	return &har.HAR{}, nil
}
