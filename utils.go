package xatago

import (
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
