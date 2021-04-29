package utils

import (
	"time"

	"github.com/konrads/go-tagged-articles/pkg/model"
)

// return unique set from the passed set, in the original insert order
func UniqueStrings(xs []string) []string {
	type void struct{}
	var member void

	keys := map[string]void{}
	res := []string{}
	for _, x := range xs {
		if _, ok := keys[x]; !ok {
			keys[x] = member
			res = append(res, x)
		}
	}
	return res
}

// filter out string in a slice of strings
func FilterStrings(xs []string, filtered string) []string {
	res := []string{}
	for _, x := range xs {
		if x != filtered {
			res = append(res, x)
		}
	}
	return res
}

// convert a date to url param format of 20060102
func ToUrlParamDate(str string) (*model.Date, error) {
	t, err := time.Parse("20060102", str)
	if err != nil {
		return nil, err
	}
	asDate := model.Date(t)
	return &asDate, nil
}
