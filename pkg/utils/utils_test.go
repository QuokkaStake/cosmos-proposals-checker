package utils

import (
	"main/pkg/constants"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func StringOfRandomLength(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

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

	require.NoError(t, err)
	assert.Equal(t, int64(0), value)
}

func TestGetBlockFromHeaderInvalidValue(t *testing.T) {
	t.Parallel()

	header := http.Header{
		constants.HeaderBlockHeight: []string{"invalid"},
	}
	value, err := GetBlockHeightFromHeader(header)

	require.Error(t, err)
	assert.Equal(t, int64(0), value)
}

func TestGetBlockFromHeaderValidValue(t *testing.T) {
	t.Parallel()

	header := http.Header{
		constants.HeaderBlockHeight: []string{"123"},
	}
	value, err := GetBlockHeightFromHeader(header)

	require.NoError(t, err)
	assert.Equal(t, int64(123), value)
}

func TestMustMarshallPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			require.Fail(t, "Expected to have a panic here!")
		}
	}()

	// see https://stackoverflow.com/a/33964549/1206421
	MustMarshal(make(chan int))
}

func TestMustMarshallValid(t *testing.T) {
	t.Parallel()

	str := "test"
	content := MustMarshal(str)
	assert.Equal(t, []byte("\"test\""), content)
}

func TestSplitStringIntoChunksLessThanOneChunk(t *testing.T) {
	t.Parallel()

	str := StringOfRandomLength(10)
	chunks := SplitStringIntoChunks(str, 20)
	assert.Len(t, chunks, 1, "There should be 1 chunk!")
}

func TestSplitStringIntoChunksExactlyOneChunk(t *testing.T) {
	t.Parallel()

	str := StringOfRandomLength(10)
	chunks := SplitStringIntoChunks(str, 10)

	assert.Len(t, chunks, 1, "There should be 1 chunk!")
}

func TestSplitStringIntoChunksMoreChunks(t *testing.T) {
	t.Parallel()

	str := "aaaa\nbbbb\ncccc\ndddd\neeeee\n"
	chunks := SplitStringIntoChunks(str, 10)
	assert.Len(t, chunks, 3, "There should be 3 chunks!")
}

func TestSubtract(t *testing.T) {
	t.Parallel()

	type TestStruct struct {
		Value string
	}

	first := []TestStruct{
		{Value: "1"},
		{Value: "2"},
		{Value: "3"},
	}

	second := []TestStruct{
		{Value: "2"},
		{Value: "4"},
	}

	result := Subtract(first, second, func(v TestStruct) any { return v.Value })
	assert.Len(t, result, 2)
	assert.Equal(t, "1", result[0].Value)
	assert.Equal(t, "3", result[1].Value)
}

func TestUnion(t *testing.T) {
	t.Parallel()

	type TestStruct struct {
		Value string
	}

	first := []TestStruct{
		{Value: "1"},
		{Value: "2"},
		{Value: "3"},
	}

	second := []TestStruct{
		{Value: "2"},
		{Value: "4"},
	}

	result := Union(first, second, func(v TestStruct) any { return v.Value })
	assert.Len(t, result, 1)
	assert.Equal(t, "2", result[0].Value)
}

func TestMapToArray(t *testing.T) {
	t.Parallel()

	testMap := map[string]string{
		"test1": "1",
		"test2": "2",
		"test3": "3",
	}

	result := MapToArray(testMap)
	assert.Len(t, result, 3)
	assert.Contains(t, result, "1")
	assert.Contains(t, result, "2")
	assert.Contains(t, result, "3")
}
