package ojson

import (
	"bytes"
	"encoding/json"
)

type Any struct {
	v any
}

func (a Any) Value() any {
	return a.v
}

func (a Any) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.v)
}

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
