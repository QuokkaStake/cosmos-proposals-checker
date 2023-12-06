package types

import "encoding/json"

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
