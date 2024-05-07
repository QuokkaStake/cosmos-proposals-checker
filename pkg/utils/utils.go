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
