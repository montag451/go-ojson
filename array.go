package ojson

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Array represents a slice of any values. When using it to decode a
// JSON array it ensures that embedded JSON objects are decoded as
// [Object].
type Array []any

// MarshalJSON implements the [json.Marshaler] interface.
func (a *Array) UnmarshalJSON(d []byte) error {
	dec := json.NewDecoder(bytes.NewReader(d))
	tok, err := dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := tok.(json.Delim); ok && delim == '[' {
		return a.unmarshalJSON(dec)
	}
	return fmt.Errorf(`expected "[", got %q`, tok)
}

// UnmarshalJSON implements the [json.Unmarshaler] interface.
func (a *Array) unmarshalJSON(d *json.Decoder) error {
	*a = (*a)[:0]
	for d.More() {
		var v Any
		if err := v.unmarshalJSON(d); err != nil {
			return err
		}
		*a = append(*a, v.Value())
	}
	if *a == nil {
		*a = []any{}
	}
	_, err := d.Token()
	return err
}
