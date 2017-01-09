package har

import (
	"github.com/fedesog/webdriver"
	"github.com/jordanpotter/har"
)

func New(logEntries []webdriver.LogEntry) (*har.HAR, error) {
	return &har.HAR{}, nil
}
