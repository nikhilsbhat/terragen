package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	sampleSlice := []string{"first", "second", "third"}
	t.Run("Should return true as the element was found in the passed slice", func(t *testing.T) {
		expected := true
		actual := Contains(sampleSlice, "first")
		assert.Equal(t, expected, actual)
	})
	t.Run("Should return false as the element was missing in the passed slice", func(t *testing.T) {
		expected := false
		actual := Contains(sampleSlice, "firstElement")
		assert.Equal(t, expected, actual)
	})
}

func TestHasChange(t *testing.T) {
	var (
		sampleSliceOne = []string{"first", "second", "third"}
	)
	t.Run("should return true as slices to be compared has changes from one another", func(t *testing.T) {
		newSampleSliceTwo := []string{"first", "second"}
		expected := true
		actual := HasChange(sampleSliceOne, newSampleSliceTwo)
		assert.Equal(t, expected, actual)
	})
	t.Run("should return false as slices to be compared has no changes from one another", func(t *testing.T) {
		sampleSliceTwo := []string{"first", "second", "third"}
		expected := false
		actual := HasChange(sampleSliceOne, sampleSliceTwo)
		assert.Equal(t, expected, actual)
	})
}
