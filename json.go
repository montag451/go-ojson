package ojson

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Any struct {
	v any
}

func (a Any) Interface() any {
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
			var o Object
			if err := o.unmarshalJSON(d); err != nil {
				return err
			}
			a.v = o
		case '[':
			var array Array
			if err := array.unmarshalJSON(d); err != nil {
				return err
			}
			a.v = array
		}
	default:
		a.v = v
	}
	return nil
}

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
	switch v := tok.(type) {
	case json.Delim:
		switch v {
		case '[':
			if err := a.unmarshalJSON(dec); err != nil {
				return fmt.Errorf("failed to unmarshal array: error at offset %d: %v", dec.InputOffset(), err)
			}
			return nil
		}
	}
	return fmt.Errorf("expected \"[\", got %q", tok)
}

func (a *Array) unmarshalJSON(d *json.Decoder) error {
	for d.More() {
		var v Any
		v.unmarshalJSON(d)
		*a = append(*a, v.Interface())
	}
	_, err := d.Token()
	return err
}

type Object struct {
	m    map[string]any
	keys []string
}

func (o *Object) Set(k string, v any) {
	if o.m == nil {
		o.m = make(map[string]any)
	}
	if _, ok := o.m[k]; !ok {
		o.keys = append(o.keys, k)
	}
	o.m[k] = v
}

func (o Object) Range(f func(k string, v any) bool) {
	for _, k := range o.keys {
		if !f(k, o.m[k]) {
			return
		}
	}
}

func (o Object) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	b.WriteString("{")
	for i, k := range o.keys {
		if i != 0 {
			b.WriteString(",")
		}
		res, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}
		b.Write(res)
		b.WriteString(":")
		res, err = json.Marshal(o.m[k])
		if err != nil {
			return nil, err
		}
		b.Write(res)
	}
	b.WriteString("}")
	return b.Bytes(), nil
}

func (o *Object) UnmarshalJSON(d []byte) error {
	dec := json.NewDecoder(bytes.NewReader(d))
	tok, err := dec.Token()
	if err != nil {
		return err
	}
	switch v := tok.(type) {
	case json.Delim:
		switch v {
		case '{':
			if err := o.unmarshalJSON(dec); err != nil {
				return fmt.Errorf("failed to unmarshal object: error at offset %d: %v", dec.InputOffset(), err)
			}
			return nil
		}
	}
	return fmt.Errorf("expected \"{\", got %q", tok)
}

func (o *Object) unmarshalJSON(d *json.Decoder) error {
	for d.More() {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		key := tok.(string)
		var v Any
		v.unmarshalJSON(d)
		o.Set(key, v.Interface())
	}
	_, err := d.Token()
	return err
}
