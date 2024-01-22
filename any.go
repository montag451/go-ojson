package ojson

import (
	"bytes"
	"encoding/json"
)

// Any represents any Go value. When using it to decode a JSON value
// it ensures that embedded JSON objects are decoded as [Object].
type Any struct {
	v any
}

// Value returns the value encapsulated in a.
func (a Any) Value() any {
	return a.v
}

// MarshalJSON implements the [json.Marshaler] interface.
func (a Any) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.v)
}

// UnmarshalJSON implements the [json.Unmarshaler] interface.
func (a *Any) UnmarshalJSON(d []byte) error {
	dec := json.NewDecoder(bytes.NewReader(d))
	return a.unmarshalJSON(dec)
}

func (a *Any) unmarshalJSON(d *json.Decoder) error {
	tok, err := d.Token()
	if err != nil {
		return err
	}
	switch v := tok.(type) {
	case json.Delim:
		switch v {
		case '{':
			o := NewObject()
			if err := o.unmarshalJSON(d); err != nil {
				return err
			}
			a.v = o
		case '[':
			var arr Array
			if err := arr.unmarshalJSON(d); err != nil {
				return err
			}
			a.v = []any(arr)
		}
	default:
		a.v = v
	}
	return nil
}
