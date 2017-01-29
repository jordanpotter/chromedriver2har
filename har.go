package chromedriver2har

import (
	"encoding/json"

	"github.com/fedesog/webdriver"
	"github.com/jordanpotter/har"
	"github.com/pkg/errors"
)

func New(logEntries []webdriver.LogEntry) (*har.HAR, error) {
	chromeLogEntries, err := chromeLogEntries(logEntries)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create chrome log entries")
	}

	_, err = harEntries(chromeLogEntries)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create HAR entries")
	}

	return &har.HAR{}, nil
}

func chromeLogEntries(logEntries []webdriver.LogEntry) ([]ChromeLogEntry, error) {
	chromeLogEntries := make([]ChromeLogEntry, 0, len(logEntries))
	for _, logEntry := range logEntries {
		var chromeEntry ChromeLogEntry
		if err := json.Unmarshal([]byte(logEntry.Message), &chromeEntry); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal log entry")
		}
		chromeLogEntries = append(chromeLogEntries, chromeEntry)
	}
	return chromeLogEntries, nil
}

func harEntries(chromeLogEntries []ChromeLogEntry) ([]har.Entry, error) {
	for _, chromeLogEntry := range chromeLogEntries {
		var err error

		switch chromeLogEntry.Message.Method {
		case "Network.requestWillBeSent":
			err = processRequestWillBeSent(chromeLogEntry.Message.Params)
		case "Network.responseReceived":
			err = processResponseReceived(chromeLogEntry.Message.Params)
		case "Network.dataReceived":
			err = processDataReceived(chromeLogEntry.Message.Params)
		case "Network.loadingFinished":
			err = processLoadingFinished(chromeLogEntry.Message.Params)
		}

		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse entry %q", chromeLogEntry.Message.Method)
		}
	}

	return nil, nil
}

func processRequestWillBeSent(params json.RawMessage) error {
	var data NetworkRequestWillBeSent
	if err := json.Unmarshal(params, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal NetworkRequestWillBeSent data")
	}
	return nil
}

func processResponseReceived(params json.RawMessage) error {
	var data NetworkResponseReceived
	if err := json.Unmarshal(params, &data); err != nil {
		errors.Wrap(err, "failed to unmarshal NetworkResponseReceived data")
	}
	return nil
}

func processDataReceived(params json.RawMessage) error {
	var data NetworkDataReceived
	if err := json.Unmarshal(params, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal NetworkDataReceived data")
	}
	return nil
}

func processLoadingFinished(params json.RawMessage) error {
	var data NetworkLoadingFinished
	if err := json.Unmarshal(params, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal NetworkLoadingFinished data")
	}
	return nil
}
