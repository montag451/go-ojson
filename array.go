package ojson

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Array []any

func (a Array) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	b.WriteString("[")
	for i, v := range a {
		if i != 0 {
			b.WriteString(",")
		}
		res, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		b.Write(res)
	}
	b.WriteString("]")
	return b.Bytes(), nil
}

func (a *Array) UnmarshalJSON(d []byte) error {
	dec := json.NewDecoder(bytes.NewReader(d))
	tok, err := dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := tok.(json.Delim); ok && delim == '[' {
		if err := a.unmarshalJSON(dec); err != nil {
			return fmt.Errorf("failed to unmarshal array: error at offset %d: %w", dec.InputOffset(), err)
		}
		return nil
	}
	return fmt.Errorf(`expected "[", got %q`, tok)
}

func (a *Array) unmarshalJSON(d *json.Decoder) error {
	for d.More() {
		var v Any
		if err := v.unmarshalJSON(d); err != nil {
			return err
		}
		*a = append(*a, v.Value())
	}
	_, err := d.Token()
	return err
}
