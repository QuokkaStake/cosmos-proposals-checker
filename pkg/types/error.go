package types

import (
	"encoding/json"
	"fmt"
	"strings"
)

type JSONError struct {
	error string
}

func NewJSONError(err error) JSONError {
	return JSONError{error: err.Error()}
}

func (e *JSONError) Error() string {
	return e.error
}

func (e *JSONError) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.error)
}

func (e *JSONError) UnmarshalJSON(data []byte) error {
	e.error = string(data)
	return nil
}

type NodeError struct {
	Node  string
	Error JSONError
}

type QueryError struct {
	QueryError error
	NodeErrors []NodeError
}

func (q QueryError) Error() string {
	if q.QueryError != nil {
		return q.QueryError.Error()
	}

	var sb strings.Builder

	sb.WriteString("All LCD requests failed:\n")
	for index, nodeError := range q.NodeErrors {
		sb.WriteString(fmt.Sprintf("#%d: %s -> %s\n", index+1, nodeError.Node, nodeError.Error.error))
	}

	return sb.String()
}
