package chromedriver2har

import (
	"encoding/json"
	"net"
)

type ChromeLogEntry struct {
	Message Message `json:"message"`
	Webview string  `json:"webview"`
}

type Message struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type NetworkRequestWillBeSent struct {
	RequestID        string    `json:"requestId"`
	LoaderID         string    `json:"loaderId"`
	DocumentURL      string    `json:"documentURL"`
	Request          Request   `json:"request"`
	Timestamp        float64   `json:"timestamp"`
	RedirectResponse *Response `json:"redirectResponse"`
}

type NetworkResponseReceived struct {
	RequestID string   `json:"requestId"`
	LoaderID  string   `json:"loaderId"`
	Timestamp float64  `json:"timestamp"`
	Type      string   `json:"type"`
	Response  Response `json:"response"`
}

type NetworkDataReceived struct {
	RequestID         string  `json:"requestId"`
	Timestamp         float64 `json:"timestamp"`
	DataLength        int     `json:"dataLength"`
	EncodedDataLength int     `json:"encodedDataLength"`
}

type NetworkLoadingFinished struct {
	RequestID         string  `json:"requestId"`
	Timestamp         float64 `json:"timestamp"`
	EncodedDataLength int     `json:"encodedDataLength"`
}

type Request struct {
	URL              string            `json:"url"`
	Method           string            `json:"method"`
	Headers          map[string]string `json:"headers"`
	PostData         *string           `json:"postData"`
	MixedContentType *string           `json:"mixedContentType"`
	InitialPriority  string            `json:"initialPriority"`
}

type Response struct {
	URL                string                 `json:"url"`
	Status             int                    `json:"status"`
	StatusText         string                 `json:"statusText"`
	Headers            map[string]string      `json:"headers"`
	HeadersText        *string                `json:"headersText"`
	MimeType           string                 `json:"mimeType"`
	RequestHeaders     map[string]string      `json:"requestHeaders"`
	RequestHeadersText *string                `json:"requestHeadersText"`
	ConnectionReused   bool                   `json:"connectionReused"`
	ConnectionID       int                    `json:"connectionId"`
	RemoteIPAddress    *net.IP                `json:"remoteIPAddress"`
	RemotePort         *int                   `json:"remotePort"`
	FromDiskCache      *bool                  `json:"fromDiskCache"`
	FromServiceWorker  *bool                  `json:"fromServiceWorker"`
	EncodedDataLength  int                    `json:"encodedDataLength"`
	Timing             *Timing                `json:"timing"`
	Protocol           *string                `json:"protocol"`
	SecurityState      string                 `json:"securityState"`
	SecurityDetails    map[string]interface{} `json:"securityDetails"`
}

type Timing struct {
	RequestTime       float64 `json:"requestTime"`
	ProxyStart        float64 `json:"proxyStart"`
	ProxyEnd          float64 `json:"proxyEnd"`
	DNSStart          float64 `json:"dnsStart"`
	DNSEnd            float64 `json:"dnsEnd"`
	ConnectStart      float64 `json:"connectStart"`
	ConnectEnd        float64 `json:"connectEnd"`
	SSLStart          float64 `json:"sslStart"`
	SSLEnd            float64 `json:"sslEnd"`
	SendStart         float64 `json:"sendStart"`
	SendEnd           float64 `json:"sendEnd"`
	ReceiveHeadersEnd float64 `json:"receiveHeadersEnd"`
}
