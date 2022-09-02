package decode_test

import (
	"fmt"
	"testing"

	"github.com/nikhilsbhat/terragen/pkg/decode"

	"github.com/stretchr/testify/assert"
)

func TestJSONDecode(t *testing.T) {
	type testStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	t.Run("should be able to unmarshal byte to json", func(t *testing.T) {
		jsonData := `{"name":"John", "age":30}`
		name := &testStruct{}

		expected := &testStruct{Age: 30, Name: "John"}

		err := decode.JSONDecode([]byte(jsonData), name)
		assert.NoError(t, err)
		assert.Equal(t, expected, name)
	})

	t.Run("should throw unknown type error", func(t *testing.T) {
		jsonData := `{"name":30, "age":30}`
		name := &testStruct{}

		expectedErr := "error occurred at line 1, json: cannot unmarshal number into Go struct field testStruct.name of type " +
			"string\n{\"name\":30, \"age\":30}\nThe data type you entered for the value is wrong"

		err := decode.JSONDecode([]byte(jsonData), name)
		assert.EqualError(t, err, expectedErr)
	})

	t.Run("should throw syntaxError", func(t *testing.T) {
		jsonData := `{"name":30, "age:30}`
		name := &testStruct{}
		expected := &testStruct{Name: "", Age: 0}

		expectedErr := "error occurred at line 1, unexpected end of JSON input\n{\"name\":30, \"age:30}"

		err := decode.JSONDecode([]byte(jsonData), name)
		assert.EqualError(t, err, expectedErr)
		assert.Equal(t, name, expected)
	})
}

func TestGetStringOfMessage(t *testing.T) {
	t.Run("should return string of error message", func(t *testing.T) {
		customErr := fmt.Errorf("this is a error message") //nolint:goerr113
		expected := "this is a error message"

		actual := decode.GetStringOfMessage(customErr)
		assert.Equal(t, expected, actual)
	})
}
