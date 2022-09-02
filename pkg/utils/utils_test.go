package utils_test

import (
	"testing"

	"github.com/nikhilsbhat/terragen/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	sampleSlice := []string{"first", "second", "third"}
	t.Run("Should return true as the element was found in the passed slice", func(t *testing.T) {
		actual := utils.Contains(sampleSlice, "first")
		assert.Equal(t, true, actual)
	})
	t.Run("Should return false as the element was missing in the passed slice", func(t *testing.T) {
		actual := utils.Contains(sampleSlice, "firstElement")
		assert.Equal(t, false, actual)
	})
}

func TestHasChange(t *testing.T) {
	sampleSliceOne := []string{"first", "second", "third"}
	t.Run("should return true as slices to be compared has changes from one another", func(t *testing.T) {
		newSampleSliceTwo := []string{"first", "second"}
		actual := utils.HasChange(sampleSliceOne, newSampleSliceTwo)
		assert.Equal(t, true, actual)
	})
	t.Run("should return false as slices to be compared has no changes from one another", func(t *testing.T) {
		sampleSliceTwo := []string{"first", "second", "third"}
		actual := utils.HasChange(sampleSliceOne, sampleSliceTwo)
		assert.Equal(t, false, actual)
	})
}
