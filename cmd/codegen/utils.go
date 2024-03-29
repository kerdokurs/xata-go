package main

import (
	"fmt"
	"os"
)

func errExit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Generation error: %v\n", err)
		os.Exit(1)
	}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func filter[T any](ts []T, f func(T) bool) []T {
	var res []T
	for _, t := range ts {
		if f(t) {
			res = append(res, t)
		}
	}
	return res
}
