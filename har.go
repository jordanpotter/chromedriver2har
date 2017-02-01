package chromedriver2har

import (
	"encoding/json"

	"github.com/fedesog/webdriver"
	"github.com/jordanpotter/har"
	"github.com/pkg/errors"
)

const (
	harVersion     = "1.2"
	creatorName    = "chromedriver2har"
	creatorVersion = "0.1"
)

func New(logEntries []webdriver.LogEntry) (*har.HAR, error) {
	chromeLogEntries, err := chromeLogEntries(logEntries)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create chrome log entries")
	}

	page, err := harPage(chromeLogEntries)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create HAR page")
	}

	entries, err := harEntries(chromeLogEntries)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create HAR entries")
	}

	return &har.HAR{
		Log: har.Log{
			Version: harVersion,
			Creator: har.Creator{
				Name:    creatorName,
				Version: creatorVersion,
			},
			Pages:   []har.Page{page},
			Entries: entries,
		},
	}, nil
}

func chromeLogEntries(logEntries []webdriver.LogEntry) ([]ChromeLogEntry, error) {
	chromeLogEntries := make([]ChromeLogEntry, 0, len(logEntries))
	for _, logEntry := range logEntries {
		var chromeEntry ChromeLogEntry
		if err := json.Unmarshal([]byte(logEntry.Message), &chromeEntry); err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal log entry at timestamp %d", logEntry.TimeStamp)
		}
		chromeLogEntries = append(chromeLogEntries, chromeEntry)
	}
	return chromeLogEntries, nil
}
