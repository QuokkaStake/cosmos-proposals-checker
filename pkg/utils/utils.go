package utils

import (
	"encoding/json"
	"fmt"
	"main/pkg/constants"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Filter[T any](slice []T, f func(T) bool) []T {
	var n []T
	for _, e := range slice {
		if f(e) {
			n = append(n, e)
		}
	}
	return n
}

func Map[T, V any](slice []T, f func(T) V) []V {
	result := make([]V, len(slice))

	for index, value := range slice {
		result[index] = f(value)
	}

	return result
}

func Contains[T comparable](slice []T, elt T) bool {
	for _, value := range slice {
		if value == elt {
			return true
		}
	}

	return false
}

// Subtract returns a new slice than includes values that are presented in the first, but not
// the second array.
func Subtract[T any, C comparable](first, second []T, predicate func(T) C) []T {
	valuesMap := make(map[C]bool, len(second))
	for _, value := range second {
		valuesMap[predicate(value)] = true
	}

	newSlice := make([]T, 0)

	for _, value := range first {
		predicateResult := predicate(value)
		_, ok := valuesMap[predicateResult]
		if !ok {
			newSlice = append(newSlice, value)
		}
	}

	return newSlice
}

func Union[T any, C comparable](first, second []T, predicate func(T) C) []T {
	valuesMap := make(map[C]bool, len(second))
	for _, value := range second {
		valuesMap[predicate(value)] = true
	}

	newSlice := make([]T, 0)

	for _, value := range first {
		predicateResult := predicate(value)
		_, ok := valuesMap[predicateResult]
		if ok {
			newSlice = append(newSlice, value)
		}
	}

	return newSlice
}

func MapToArray[K comparable, T any](source map[K]T) []T {
	newSlice := make([]T, len(source))

	index := 0

	for _, value := range source {
		newSlice[index] = value
		index++
	}

	return newSlice
}

func FormatDuration(duration time.Duration) string {
	days := int64(duration.Hours() / 24)
	hours := int64(math.Mod(duration.Hours(), 24))
	minutes := int64(math.Mod(duration.Minutes(), 60))
	seconds := int64(math.Mod(duration.Seconds(), 60))

	chunks := []struct {
		singularName string
		amount       int64
	}{
		{"day", days},
		{"hour", hours},
		{"minute", minutes},
		{"second", seconds},
	}

	parts := []string{}

	for _, chunk := range chunks {
		switch chunk.amount {
		case 0:
			continue
		case 1:
			parts = append(parts, fmt.Sprintf("%d %s", chunk.amount, chunk.singularName))
		default:
			parts = append(parts, fmt.Sprintf("%d %ss", chunk.amount, chunk.singularName))
		}
	}

	return strings.Join(parts, " ")
}

func GetBlockHeightFromHeader(header http.Header) (int64, error) {
	valueStr := header.Get(constants.HeaderBlockHeight)
	if valueStr == "" {
		return 0, nil
	}

	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func MustMarshal(v any) []byte {
	if content, err := json.Marshal(v); err != nil {
		panic(err)
	} else {
		return content
	}
}

func SplitStringIntoChunks(msg string, maxLineLength int) []string {
	msgsByNewline := strings.Split(msg, "\n")
	outMessages := []string{}

	var sb strings.Builder

	for _, line := range msgsByNewline {
		if sb.Len()+len(line) > maxLineLength {
			outMessages = append(outMessages, sb.String())
			sb.Reset()
		}

		sb.WriteString(line + "\n")
	}

	outMessages = append(outMessages, sb.String())
	return outMessages
}
