package utils

import (
	"bytes"
	"encoding/gob"
)

// Contains returns true if given element is present the specified slice.
func Contains(s []string, searchTerm string) bool {
	for _, i := range s {
		if i == searchTerm {
			return true
		}
	}

	return false
}

// HasChange true if two slices has differences.
func HasChange(slice1, slice2 []string) bool {
	return !bytes.Equal(byts(slice1), byts(slice2))
}

func byts(slice []string) []byte {
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(slice) //nolint:errcheck

	return buf.Bytes()
}
