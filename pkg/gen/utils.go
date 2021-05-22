package gen

import (
	"bytes"
	"encoding/gob"
)

// contains returns true if given element is present the specified slice.
func contains(s []string, searchTerm string) bool {
	for _, i := range s {
		if i == searchTerm {
			return true
		}
	}
	return false
}

func hasChange(slice1, slice2 []string) bool {
	return !bytes.Equal(byts(slice1), byts(slice2))
}

func byts(slice []string) []byte {
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(slice) //nolint:errcheck
	return buf.Bytes()
}
