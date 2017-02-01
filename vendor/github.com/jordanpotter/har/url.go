package har

import "net/url"

type URL struct {
	url.URL
}

func (u URL) MarshalJSON() ([]byte, error) {
	return []byte(`"` + u.String() + `"`), nil
}
