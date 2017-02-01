package chromedriver2har

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/jordanpotter/har"
	"github.com/pkg/errors"
)

func harEntries(chromeLogEntries []ChromeLogEntry) ([]har.Entry, error) {
	paramsByRequest, err := paramsByRequest(chromeLogEntries)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create params by request")
	}

	entries := make([]har.Entry, 0, len(paramsByRequest))
	for requestID, params := range paramsByRequest {
		if !params.completed() {
			continue
		}

		entry, err := harEntry(params)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create har entry for request %q", requestID)
		}

		entries = append(entries, entry)
	}

	// TODO: sort entries
	return entries, nil
}

func harEntry(params *requestParams) (har.Entry, error) {
	request, err := harRequest(params)
	if err != nil {
		return har.Entry{}, errors.Wrap(err, "failed to create har request")
	}

	response, err := harResponse(params)
	if err != nil {
		return har.Entry{}, errors.Wrap(err, "failed to create har response")
	}

	cache, err := harCache(params)
	if err != nil {
		return har.Entry{}, errors.Wrap(err, "failed to create har cache")
	}

	timings, err := harTimings(params)
	if err != nil {
		return har.Entry{}, errors.Wrap(err, "failed to create har timinhgs")
	}

	return har.Entry{
		StartedDateTime: harEntryStartedDateTime(params),
		Time:            harEntryTime(params),
		Request:         request,
		Response:        response,
		Cache:           cache,
		Timings:         timings,
	}, nil
}

func harEntryStartedDateTime(params *requestParams) har.Time {
	wallTimeNanoseconds := params.networkRequestWillBeSent.WallTime * float64(time.Second) / float64(time.Nanosecond)
	startedDateTime := time.Unix(0, int64(wallTimeNanoseconds))
	return har.Time{startedDateTime}
}

func harEntryTime(params *requestParams) float64 {
	start := params.networkRequestWillBeSent.Timestamp
	end := params.networkLoadingFinished.Timestamp
	return (end - start) * 1000
}

func harRequest(params *requestParams) (har.Request, error) {
	request := params.networkRequestWillBeSent.Request
	response := params.networkResponseReceived.Response

	requestURL, err := url.Parse(request.URL)
	if err != nil {
		return har.Request{}, errors.Wrapf(err, "failed to parse url %q", request.URL)
	}

	bodySize := -1
	if request.PostData != nil {
		bodySize = len(*request.PostData)
	}

	headersSize := harRequestHeadersSize(request.Method, safeStringDereference(response.Protocol), *requestURL, request.Headers)

	return har.Request{
		Method:      request.Method,
		URL:         har.URL{*requestURL},
		HTTPVersion: safeStringDereference(response.Protocol),
		Cookies:     harCookies(request.Headers),
		Headers:     harHeaders(request.Headers),
		QueryString: harQueryStringParams(*requestURL),
		// PostData: TODO,
		HeadersSize: headersSize,
		BodySize:    bodySize,
	}, nil
}

func harResponse(params *requestParams) (har.Response, error) {
	response := params.networkResponseReceived.Response

	redirectURL := &url.URL{}
	if params.networkRequestWillBeSentRedirect != nil {
		var err error
		redirectResponse := params.networkRequestWillBeSentRedirect.RedirectResponse
		redirectURL, err = url.Parse(redirectResponse.URL)
		if err != nil {
			return har.Response{}, errors.Wrapf(err, "failed to parse url %q", redirectResponse.URL)
		}
	}

	headersSize := harResponseHeadersSize(safeStringDereference(response.Protocol), response.Status, response.StatusText, response.Headers)
	bodySize := params.networkLoadingFinished.EncodedDataLength - headersSize

	return har.Response{
		Status:      response.Status,
		StatusText:  response.StatusText,
		HTTPVersion: safeStringDereference(response.Protocol),
		Cookies:     harCookies(response.Headers),
		Headers:     harHeaders(response.Headers),
		Content:     harContent(params, bodySize),
		RedirectURL: har.URL{*redirectURL},
		HeadersSize: headersSize,
		BodySize:    bodySize,
	}, nil
}

func harCache(params *requestParams) (har.Cache, error) {
	return har.Cache{}, nil
}

func harTimings(params *requestParams) (har.Timings, error) {
	response := params.networkResponseReceived.Response

	if response.Timing == nil {
		return har.Timings{}, nil
	}

	correctTiming := func(num float64) float64 {
		if num == 0.0 {
			return -1
		}
		return num
	}

	blocked := correctTiming(response.Timing.DNSStart)
	dns := correctTiming(response.Timing.DNSEnd - response.Timing.DNSStart)
	connect := correctTiming(response.Timing.ConnectEnd - response.Timing.ConnectStart)
	send := correctTiming(response.Timing.SendEnd - response.Timing.SendStart)
	wait := correctTiming(response.Timing.ReceiveHeadersEnd - response.Timing.SendEnd)
	receive := params.networkLoadingFinished.Timestamp*1000 - params.networkRequestWillBeSent.Timestamp*1000 - response.Timing.ReceiveHeadersEnd
	ssl := correctTiming(response.Timing.SSLEnd - response.Timing.SSLStart)

	return har.Timings{
		Blocked: &blocked,
		DNS:     &dns,
		Connect: &connect,
		Send:    send,
		Wait:    wait,
		Receive: receive,
		SSL:     &ssl,
	}, nil
}

func harCookies(headers map[string]string) []har.Cookie {
	harCookies := make([]har.Cookie, 0)
	for key, value := range headers {
		if key == "Cookie" {
			cookieStrs := strings.Split(value, ";")
			for _, cookieStr := range cookieStrs {
				harCookie := harCookie(cookieStr)
				harCookies = append(harCookies, harCookie)
			}
		}
	}
	return harCookies
}

func harCookie(cookie string) har.Cookie {
	cookie = strings.TrimSpace(cookie)
	components := strings.SplitN(cookie, "=", 2)
	return har.Cookie{Name: components[0], Value: components[1]}
}

func harHeaders(headers map[string]string) []har.Header {
	harHeaders := make([]har.Header, 0, len(headers))
	for key, value := range headers {
		harHeader := har.Header{Name: key, Value: value}
		harHeaders = append(harHeaders, harHeader)
	}
	return harHeaders
}

func harRequestHeadersSize(method, protocol string, u url.URL, headers map[string]string) int {
	size := len(fmt.Sprintf("%s %s %s\r\n", method, u.RequestURI(), protocol))
	for key, value := range headers {
		size += len(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	return size + len("\r\n")
}

func harResponseHeadersSize(protocol string, statusCode int, statusText string, headers map[string]string) int {
	size := len(fmt.Sprintf("%s %d %s\r\n", protocol, statusCode, statusText))
	for key, value := range headers {
		size += len(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	return size + len("\r\n")
}

func harQueryStringParams(u url.URL) []har.QueryStringParam {
	harQueryStringParams := make([]har.QueryStringParam, 0)
	for key, values := range u.Query() {
		for _, value := range values {
			harQueryStringParam := har.QueryStringParam{Name: key, Value: value}
			harQueryStringParams = append(harQueryStringParams, harQueryStringParam)
		}
	}
	return harQueryStringParams
}

func harContent(params *requestParams, bodySize int) har.Content {
	response := params.networkResponseReceived.Response

	size := 0
	for _, dataReceived := range params.networkDatasReceived {
		size += dataReceived.DataLength
	}

	compression := size - bodySize

	return har.Content{
		Size:        size,
		Compression: &compression,
		MIMEType:    response.MimeType,
		// Text: TODO,
		// Encoding: TODO,
	}
}
