package chromedriver2har

import (
	"encoding/json"

	"github.com/pkg/errors"
)

const (
	MethodNetworkRequestWillBeSent = "Network.requestWillBeSent"
	MethodNetworkResponseReceived  = "Network.responseReceived"
	MethodNetworkDataReceived      = "Network.dataReceived"
	MethodNetworkLoadingFinished   = "Network.loadingFinished"
)

type requestParams struct {
	networkRequestWillBeSent         NetworkRequestWillBeSent
	networkRequestWillBeSentRedirect *NetworkRequestWillBeSent
	networkResponseReceived          NetworkResponseReceived
	networkDatasReceived             []NetworkDataReceived
	networkLoadingFinished           NetworkLoadingFinished
}

func (rp *requestParams) completed() bool {
	return rp.networkLoadingFinished.RequestID != ""
}

func paramsByRequest(chromeLogEntries []ChromeLogEntry) (map[string]*requestParams, error) {
	paramsByRequest := make(map[string]*requestParams)

	for _, chromeLogEntry := range chromeLogEntries {
		var err error

		switch chromeLogEntry.Message.Method {
		case MethodNetworkRequestWillBeSent:
			err = processNetworkRequestWillBeSent(paramsByRequest, chromeLogEntry.Message.Params)
		case MethodNetworkResponseReceived:
			err = processNetworkResponseReceived(paramsByRequest, chromeLogEntry.Message.Params)
		case MethodNetworkDataReceived:
			err = processNetworkDataReceived(paramsByRequest, chromeLogEntry.Message.Params)
		case MethodNetworkLoadingFinished:
			err = processNetworkLoadingFinished(paramsByRequest, chromeLogEntry.Message.Params)
		}

		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse entry %q", chromeLogEntry.Message.Method)
		}
	}

	return paramsByRequest, nil
}

func processNetworkRequestWillBeSent(paramsByRequest map[string]*requestParams, params json.RawMessage) error {
	var data NetworkRequestWillBeSent
	if err := json.Unmarshal(params, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal NetworkRequestWillBeSent data")
	}

	if data.RedirectResponse != nil {
		return processNetworkRequestWillBeSentRedirect(paramsByRequest, params)
	}

	if _, ok := paramsByRequest[data.RequestID]; ok {
		return errors.Errorf("entry already exists for request %q", data.RequestID)
	}

	paramsByRequest[data.RequestID] = &requestParams{networkRequestWillBeSent: data}
	return nil
}

func processNetworkRequestWillBeSentRedirect(paramsByRequest map[string]*requestParams, params json.RawMessage) error {
	var data NetworkRequestWillBeSent
	if err := json.Unmarshal(params, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal NetworkRequestWillBeSent data")
	}

	if data.RedirectResponse == nil {
		return errors.Errorf("missing redirect response for request %q", data.RequestID)
	}

	request, ok := paramsByRequest[data.RequestID]
	if !ok {
		return errors.Errorf("missing entry for request %q", data.RequestID)
	}

	request.networkRequestWillBeSentRedirect = &data
	return nil
}

func processNetworkResponseReceived(paramsByRequest map[string]*requestParams, params json.RawMessage) error {
	var data NetworkResponseReceived
	if err := json.Unmarshal(params, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal NetworkResponseReceived data")
	}

	request, ok := paramsByRequest[data.RequestID]
	if !ok {
		return errors.Errorf("missing entry for request %q", data.RequestID)
	}

	request.networkResponseReceived = data
	return nil
}

func processNetworkDataReceived(paramsByRequest map[string]*requestParams, params json.RawMessage) error {
	var data NetworkDataReceived
	if err := json.Unmarshal(params, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal NetworkDataReceived data")
	}

	request, ok := paramsByRequest[data.RequestID]
	if !ok {
		return errors.Errorf("missing entry for request %q", data.RequestID)
	}

	request.networkDatasReceived = append(request.networkDatasReceived, data)
	return nil
}

func processNetworkLoadingFinished(paramsByRequest map[string]*requestParams, params json.RawMessage) error {
	var data NetworkLoadingFinished
	if err := json.Unmarshal(params, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal NetworkLoadingFinished data")
	}

	request, ok := paramsByRequest[data.RequestID]
	if !ok {
		return errors.Errorf("missing entry for request %q", data.RequestID)
	}

	request.networkLoadingFinished = data
	return nil
}
