package utils

import (
	"testing"
	"time"

	"github.com/konrads/go-tagged-articles/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestUniqueStrings(t *testing.T) {
	in := []string{"hello", "ball", "hi", "hi", "ball", "hi"}
	exp := []string{"hello", "ball", "hi"}
	assert.Equal(t, exp, UniqueStrings(in))
}

func TestFilter(t *testing.T) {
	in := []string{"hello", "ball", "hi", "hi", "ball", "hi"}
	exp := []string{"hello", "ball", "ball"}
	assert.Equal(t, exp, FilterStrings(in, "hi"))
}

func TestToUrlParamDate(t *testing.T) {
	in := "20170702"
	exp := time.Date(2017, 7, 2, 0, 0, 0, 0, time.UTC)
	expDate := model.Date(exp)
	res, err := ToUrlParamDate(in)
	assert.NoError(t, err)
	assert.Equal(t, &expDate, res)
}
