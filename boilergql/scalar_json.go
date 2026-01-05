package boilergql

import (
	"encoding/json"
	"io"
)

// JSON can hold any valid JSON value (object, array, string, number, boolean, null)
type JSON json.RawMessage

func (b *JSON) UnmarshalGQL(v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	*b = bytes
	return nil
}

func (b JSON) MarshalGQL(w io.Writer) {
	if len(b) == 0 {
		_, _ = w.Write([]byte("null"))
		return
	}
	_, _ = w.Write(b)
}

// UnmarshalJSON implements json.Unmarshaler
func (b *JSON) UnmarshalJSON(data []byte) error {
	*b = append((*b)[0:0], data...)
	return nil
}

// MarshalJSON implements json.Marshaler
func (b JSON) MarshalJSON() ([]byte, error) {
	if len(b) == 0 {
		return []byte("null"), nil
	}
	return b, nil
}
