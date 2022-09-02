// Package decode can decode the json to the required string.
package decode

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// JSONDecode decodes the data to json sent to it.
func JSONDecode(data []byte, i interface{}) error {
	if err := json.Unmarshal(data, i); err != nil {
		switch err.(type) { //nolint:errorlint
		case *json.UnmarshalTypeError: //nolint:errorlint
			return unknownTypeError(data, err)
		case *json.SyntaxError: //nolint:errorlint
			return syntaxError(data, err)
		}
	}

	return nil
}

func syntaxError(data []byte, err error) error {
	syntaxErr, ok := err.(*json.SyntaxError) //nolint:errorlint
	if !ok {
		return err
	}

	newline := []byte{'\x0a'}

	start := bytes.LastIndex(data[:syntaxErr.Offset], newline) + 1
	end := len(data)
	if index := bytes.Index(data[start:], newline); index >= 0 {
		end = start + index
	}

	line := bytes.Count(data[:start], newline) + 1

	err = fmt.Errorf("error occurred at line %d, %w\n%s",
		line, syntaxErr, data[start:end])

	return err
}

func unknownTypeError(data []byte, err error) error {
	unknownTypeErr, ok := err.(*json.UnmarshalTypeError) //nolint:errorlint
	if !ok {
		return err
	}

	newline := []byte{'\x0a'}

	start := bytes.LastIndex(data[:unknownTypeErr.Offset], newline) + 1
	end := len(data)
	if index := bytes.Index(data[start:], newline); index >= 0 {
		end = start + index
	}

	line := bytes.Count(data[:start], newline) + 1

	err = fmt.Errorf("error occurred at line %d, %w\n%s\nThe data type you entered for the value is wrong",
		line, unknownTypeErr, data[start:end])

	return err
}

// GetStringOfMessage returns string form of error.
func GetStringOfMessage(g interface{}) string {
	switch typ := g.(type) {
	case string:
		return typ
	case error:
		return typ.Error()
	default:
		return "unknown messagetype"
	}
}
