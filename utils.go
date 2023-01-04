package xatago

import (
	"bytes"
	"encoding/json"
	"github.com/mitchellh/mapstructure"
)

func MapToStruct[T, D any](raw T, out *D) error {
	return mapstructure.Decode(raw, out)
}

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

func MapToStructs[T, D any](ts []T) ([]D, error) {
	ds := make([]D, len(ts))

	var err error
	for i, t := range ts {
		if err = MapToStruct[T, D](t, &ds[i]); err != nil {
			return nil, err
		}
	}

	return ds, nil
}

func preprocessForCreate[T any](item *T) (map[string]any, error) {
	preprocessed := make(map[string]any)

	buf := new(bytes.Buffer)

	// Omit empty fields
	if err := json.NewEncoder(buf).Encode(item); err != nil {
		return nil, err
	}
	if err := json.NewDecoder(buf).Decode(&preprocessed); err != nil {
		return nil, err
	}

	// Remove nested fields and replace with ID if exists in the nested map
	for key, value := range preprocessed {
		if mapValue, ok := value.(map[string]any); ok {
			if idValue, ok := mapValue["id"]; ok {
				preprocessed[key] = idValue
			} else {
				delete(preprocessed, key)
			}
		}
	}

	return preprocessed, nil
}
