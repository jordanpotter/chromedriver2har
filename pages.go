package chromedriver2har

import (
	"encoding/json"

	"github.com/jordanpotter/har"
	"github.com/pkg/errors"
)

const (
	MethodPageDOMContentEventFired = "Page.domContentEventFired"
	MethodPageLoadEventFired       = "Page.loadEventFired"
)

type pageParams struct {
	firstNetworkRequestWillBeSent NetworkRequestWillBeSent
	pageDOMContentEventFired      PageDOMContentEventFired
	pageLoadEventFired            PageLoadEventFired
}

func harPage(chromeLogEntries []ChromeLogEntry) (har.Page, error) {
	params, err := pageParams(chromeLogEntries)
	if err != nil {
		return har.Page{}, errors.Wrap(err, "failed to parse page params")
	}

	return har.Page{
		//			StartedDateTime: TODO,
		ID:    "page_1",
		Title: params.firstNetworkRequestWillBeSent.DocumentURL,
		PageTimings: har.PageTimings{
			OnContentLoad: params.pageDOMContentEventFired.Timestamp,
			OnLoad:        params.pageLoadEventFired.Timestamp,
		},
	}, nil
}

func pageParams(chromeLogEntries []ChromeLogEntry) (*pageParams, error) {
	params := pageParams{}

	for _, chromeLogEntry := range chromeLogEntries {
		var err error

		switch chromeLogEntry.Message.Method {
		case MethodNetworkRequestWillBeSent:
			// TODO
		case MethodPageDOMContentEventFired:
			err = processPageDOMContentEventFired(&params, chromeLogEntry.Message.Params)
		case MethodPageLoadEventFired:
			err = processPageLoadEventFired(&params, chromeLogEntry.Message.Params)
		}

		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse entry %q", chromeLogEntry.Message.Method)
		}
	}

	return params, nil
}

func processPageDOMContentEventFired(params *pageParams, params json.RawMessage) error {
	var data PageDOMContentEventFired
	if err := json.Unmarshal(params, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal PageDOMContentEventFired data")
	}

	if params.pageDOMContentEventFired.Timestamp != 0 {
		return errors.New("already processed PageDOMContentEventFired")
	}

	params.pageDOMContentEventFired = data
	return nil
}

func processPageLoadEventFired(params *pageParams, params json.RawMessage) error {
	var data PageLoadEventFired
	if err := json.Unmarshal(params, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal PageLoadEventFired data")
	}

	if params.pageLoadEventFired.Timestamp != 0 {
		return errors.New("already processed PageLoadEventFired")
	}

	params.pageLoadEventFired = data
	return nil
}
