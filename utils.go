package xatago

import (
	"bytes"
	"encoding/json"
)

// MapToStruct converts any given map to a given generic type struct
func MapToStruct[T any](raw any) (*T, error) {
	d, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}
	bb := bytes.NewBuffer(d)
	t := new(T)
	return t, json.NewDecoder(bb).Decode(t)
}

// Map maps elements in given slice by the mapper function
func Map[T, D any](ts []T, mapper func(t T) (D, error)) ([]D, error) {
	ds := make([]D, len(ts))

	var err error
	for i, t := range ts {
		if ds[i], err = mapper(t); err != nil {
			return nil, err
		}
	}

	return ds, err
}
