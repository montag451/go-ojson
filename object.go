package ojson

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type objectValue struct {
	i int
	v any
}

type Object struct {
	m    map[string]objectValue
	keys []string
}

func (o *Object) Set(k string, v any) {
	if o.m == nil {
		o.m = make(map[string]objectValue)
	}
	oval, ok := o.m[k]
	if !ok {
		oval.i = len(o.keys)
		o.keys = append(o.keys, k)
	}
	oval.v = v
	o.m[k] = oval
}

func (o Object) Get(k string) (any, bool) {
	v, ok := o.m[k]
	return v.v, ok
}

func (o *Object) Delete(k string) {
	v, ok := o.m[k]
	if !ok {
		return
	}
	delete(o.m, k)
	o.keys = append(o.keys[:v.i], o.keys[v.i+1:]...)
}

func (o Object) Range(f func(k string, v any) bool) {
	for _, k := range o.keys {
		if !f(k, o.m[k].v) {
			return
		}
	}
}

func (o Object) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	b.WriteRune('{')
	for i, k := range o.keys {
		if i != 0 {
			b.WriteRune(',')
		}
		res, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}
		b.Write(res)
		b.WriteRune(':')
		res, err = json.Marshal(o.m[k].v)
		if err != nil {
			return nil, err
		}
		b.Write(res)
	}
	b.WriteRune('}')
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
