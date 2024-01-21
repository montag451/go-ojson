package ojson

import (
	"bytes"
	"encoding/json"
	"fmt"
)

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

func (o Object) Get(k string) (any, bool) {
	v, ok := o.m[k]
	return v, ok
}

func (o *Object) Delete(k string) {
	delete(o.m, k)
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
	if delim, ok := tok.(json.Delim); ok && delim == '{' {
		if err := o.unmarshalJSON(dec); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf(`expected "{", got %q`, tok)
}

func (o *Object) unmarshalJSON(d *json.Decoder) error {
	for d.More() {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		key := tok.(string)
		var v Any
		if err := v.unmarshalJSON(d); err != nil {
			return err
		}
		o.Set(key, v.Value())
	}
	_, err := d.Token()
	return err
}
