package decode

import (
	"fmt"
	"testing"

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

		err := JSONDecode([]byte(jsonData), name)
		assert.NoError(t, err)
		assert.Equal(t, expected, name)
	})

	t.Run("should throw unknown type error", func(t *testing.T) {
		jsonData := `{"name":30, "age":30}`
		name := &testStruct{}

		expectedErr := "error occurred at line 1, json: cannot unmarshal number into Go struct field testStruct.name of type " +
			"string\n{\"name\":30, \"age\":30}\nThe data type you entered for the value is wrong"

		err := JSONDecode([]byte(jsonData), name)
		assert.EqualError(t, err, expectedErr)
	})

	t.Run("should throw syntaxError", func(t *testing.T) {
		jsonData := `{"name":30, "age:30}`
		name := &testStruct{}
		expected := &testStruct{Name: "", Age: 0}

		expectedErr := "error occurred at line 1, unexpected end of JSON input\n{\"name\":30, \"age:30}"

		err := JSONDecode([]byte(jsonData), name)
		assert.EqualError(t, err, expectedErr)
		assert.Equal(t, name, expected)
	})
}

func TestGetStringOfMessage(t *testing.T) {
	t.Run("should return string of error message", func(t *testing.T) {
		customErr := fmt.Errorf("this is a error message")
		expected := "this is a error message"

		actual := GetStringOfMessage(customErr)
		assert.Equal(t, expected, actual)
	})
}
