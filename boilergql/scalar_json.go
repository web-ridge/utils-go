package boilergql

import (
	"encoding/json"
	"fmt"
	"io"
)

type JSON map[string]interface{}

func (b *JSON) UnmarshalGQL(v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, b)
}

func (b JSON) MarshalGQL(w io.Writer) {
	bytes, err := json.Marshal(b)
	if err != nil {
		fmt.Fprintf(w, `"%s"`, err.Error())
		return
	}
	_, _ = w.Write(bytes)
}
