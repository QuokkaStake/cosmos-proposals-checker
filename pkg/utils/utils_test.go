package utils

import (
	"github.com/stretchr/testify/assert"
	"main/pkg/constants"
	"net/http"
	"testing"
	"time"
)

func TestFilter(t *testing.T) {
	t.Parallel()

	array := []int64{1, 2, 3}

	filtered := Filter(array, func(value int64) bool {
		return value == 2
	})

	assert.Len(t, filtered, 1)
	assert.Equal(t, int64(2), filtered[0])
}

func TestMap(t *testing.T) {
	t.Parallel()

	array := []int64{1, 2, 3}

	filtered := Map(array, func(value int64) int64 {
		return value * 2
	})

	assert.Len(t, filtered, 3)
	assert.Equal(t, int64(2), filtered[0])
	assert.Equal(t, int64(4), filtered[1])
	assert.Equal(t, int64(6), filtered[2])
}

func TestContains(t *testing.T) {
	t.Parallel()

	array := []int64{1, 2, 3}

	assert.True(t, Contains(array, 3))
	assert.False(t, Contains(array, 4))
}

func TestFormatDuration(t *testing.T) {
	t.Parallel()

	duration := time.Hour*24 + time.Hour*2 + time.Second*4
	formatted := FormatDuration(duration)

	assert.Equal(t, "1 day 2 hours 4 seconds", formatted)
}

func TestGetBlockFromHeaderNoValue(t *testing.T) {
	t.Parallel()

	header := http.Header{}
	value, err := GetBlockHeightFromHeader(header)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), value)
}

func TestGetBlockFromHeaderInvalidValue(t *testing.T) {
	t.Parallel()

	header := http.Header{
		constants.HeaderBlockHeight: []string{"invalid"},
	}
	value, err := GetBlockHeightFromHeader(header)

	assert.Error(t, err)
	assert.Equal(t, int64(0), value)
}

func TestGetBlockFromHeaderValidValue(t *testing.T) {
	t.Parallel()

	header := http.Header{
		constants.HeaderBlockHeight: []string{"123"},
	}
	value, err := GetBlockHeightFromHeader(header)

	assert.NoError(t, err)
	assert.Equal(t, int64(123), value)
}
