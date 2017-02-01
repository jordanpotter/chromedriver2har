package har

import (
	"net"
)

type HAR struct {
	Log Log `json:"log"`
}

type Log struct {
	Version string   `json:"version"`
	Creator Creator  `json:"creator"`
	Browser *Browser `json:"browser,omitempty"`
	Pages   []Page   `json:"pages,omitempty"`
	Entries []Entry  `json:"entries"`
	Comment *string  `json:"comment,omitempty"`
}

type Creator struct {
	Name    string  `json:"name"`
	Version string  `json:"version"`
	Comment *string `json:"comment,omitempty"`
}

type Browser struct {
	Name    string  `json:"name"`
	Version string  `json:"version"`
	Comment *string `json:"comment,omitempty"`
}

type Page struct {
	StartedDateTime Time        `json:"startedDateTime"`
	ID              string      `json:"id"`
	Title           string      `json:"title"`
	PageTimings     PageTimings `json:"pageTimings"`
	Comment         *string     `json:"comment,omitempty"`
}

type PageTimings struct {
	OnContentLoad *float64 `json:"onContentLoad,omitempty"`
	OnLoad        *float64 `json:"onLoad,omitempty"`
	Comment       *string  `json:"comment,omitempty"`
}

type Entry struct {
	PageRef         *string  `json:"pageref,omitempty"`
	StartedDateTime Time     `json:"startedDateTime"`
	Time            float64  `json:"time"`
	Request         Request  `json:"request"`
	Response        Response `json:"response"`
	Cache           Cache    `json:"cache"`
	Timings         Timings  `json:"timings"`
	ServerIPAddress *net.IP  `json:"serverIPAddress,omitempty"`
	Connection      *string  `json:"connection,omitempty"`
	Comment         *string  `json:"comment,omitempty"`
}

type Request struct {
	Method      string             `json:"method"`
	URL         URL                `json:"url"`
	HTTPVersion string             `json:"httpVersion"`
	Cookies     []Cookie           `json:"cookies"`
	Headers     []Header           `json:"headers"`
	QueryString []QueryStringParam `json:"queryString"`
	PostData    *PostData          `json:"postData,omitempty"`
	HeadersSize int                `json:"headersSize"`
	BodySize    int                `json:"bodySize"`
	Comment     *string            `json:"comment,omitempty"`
}

type Response struct {
	Status      int      `json:"status"`
	StatusText  string   `json:"statusText"`
	HTTPVersion string   `json:"httpVersion"`
	Cookies     []Cookie `json:"cookies"`
	Headers     []Header `json:"headers"`
	Content     Content  `json:"content"`
	RedirectURL URL      `json:"redirectURL"`
	HeadersSize int      `json:"headersSize"`
	BodySize    int      `json:"bodySize"`
	Comment     *string  `json:"comment,omitempty"`
}

type Cookie struct {
	Name     string  `json:"name"`
	Value    string  `json:"value"`
	Path     *string `json:"path,omitempty"`
	Domain   *string `json:"domain,omitempty"`
	Expires  *Time   `json:"expires,omitempty"`
	HTTPOnly *bool   `json:"httpOnly,omitempty"`
	Secure   *bool   `json:"secure,omitempty"`
	Comment  *string `json:"comment,omitempty"`
}

type Header struct {
	Name    string  `json:"name"`
	Value   string  `json:"value"`
	Comment *string `json:"comment,omitempty"`
}

type QueryStringParam struct {
	Name    string  `json:"name"`
	Value   string  `json:"value"`
	Comment *string `json:"comment,omitempty"`
}

type PostData struct {
	MIMEType string          `json:"mimeType"`
	Params   []PostDataParam `json:"params"`
	Text     string          `json:"text"`
	Comment  *string         `json:"comment,omitempty"`
}

type PostDataParam struct {
	Name        string  `json:"name"`
	Value       *string `json:"value,omitempty"`
	Filename    *string `json:"fileName,omitempty"`
	ContentType *string `json:"contentType,omitempty"`
	Comment     *string `json:"comment,omitempty"`
}

type Content struct {
	Size        int     `json:"size"`
	Compression *int    `json:"compression,omitempty"`
	MIMEType    string  `json:"mimeType"`
	Text        *string `json:"text,omitempty"`
	Encoding    *string `json:"encoding,omitempty"`
	Comment     *string `json:"comment,omitempty"`
}

type Cache struct {
	BeforeRequest *CacheRequest `json:"beforeRequest,omitempty"`
	AfterRequest  *CacheRequest `json:"afterRequest,omitempty"`
	Comment       *string       `json:"comment,omitempty"`
}

type CacheRequest struct {
	Expires    *Time   `json:"expires,omitempty"`
	LastAccess Time    `json:"lastAccess"`
	ETag       string  `json:"eTag"`
	HitCount   int     `json:"hitCount"`
	Comment    *string `json:"comment,omitempty"`
}

type Timings struct {
	Blocked *float64 `json:"blocked,omitempty"`
	DNS     *float64 `json:"dns,omitempty"`
	Connect *float64 `json:"connect,omitempty"`
	Send    float64  `json:"send"`
	Wait    float64  `json:"wait"`
	Receive float64  `json:"receive"`
	SSL     *float64 `json:"ssl,omitempty"`
	Comment *string  `json:"comment,omitempty"`
}
