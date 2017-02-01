package har

import "time"

type Time struct {
	time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Format(time.RFC3339) + `"`), nil
}
