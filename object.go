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

// Object represents a JSON object. It is the equivalent of
// map[string]any but when used to decode/encode a JSON object it
// preserves the object keys order. It can also be used as an ordered
// map.
type Object struct {
	m    map[string]objectValue
	keys []string
}

// NewObject creates a new [Object] ready to be used.
func NewObject() *Object {
	return &Object{
		m: make(map[string]objectValue),
	}
}

// Set sets the value for a key
func (o *Object) Set(key string, value any) {
	oval, ok := o.m[key]
	if !ok {
		oval.i = len(o.keys)
		o.keys = append(o.keys, key)
	}
	oval.v = value
	o.m[key] = oval
}

// Get returns the value stored in the object for a key, or nil if no
// value is present. The ok result indicates whether value was found
// in the object.
func (o *Object) Get(key string) (v any, ok bool) {
	ov, ok := o.m[key]
	return ov.v, ok
}

// Delete deletes the value for a key.
func (o *Object) Delete(key string) {
	v, ok := o.m[key]
	if !ok {
		return
	}
	delete(o.m, key)
	o.keys = append(o.keys[:v.i], o.keys[v.i+1:]...)
}

// Range calls f sequentially for each key and value present in the
// object. If f returns false, range stops the iteration.
func (o *Object) Range(f func(key string, value any) bool) {
	for _, k := range o.keys {
		if !f(k, o.m[k].v) {
			return
		}
	}
}

// Len returns the number of elements in the object.
func (o *Object) Len() int {
	return len(o.keys)
}

// Clear deletes all entries in the object, resulting in an empty
// object.
func (o *Object) Clear() {
	clear(o.m)
	o.keys = o.keys[:0]
}

// MarshalJSON implements the [json.Marshaler] interface.
func (o *Object) MarshalJSON() ([]byte, error) {
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

// UnmarshalJSON implements the [json.Unmarshaler] interface.
func (o *Object) UnmarshalJSON(d []byte) error {
	dec := json.NewDecoder(bytes.NewReader(d))
	tok, err := dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := tok.(json.Delim); ok && delim == '{' {
		return o.unmarshalJSON(dec)
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
