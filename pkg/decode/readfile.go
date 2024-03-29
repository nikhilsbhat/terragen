package decode

import (
	"os"
)

// ReadFile as name specifies it just reads the content of file and returns it.
func ReadFile(filename string) ([]byte, error) {
	if _, dirneuerr := os.Stat(filename); os.IsNotExist(dirneuerr) {
		return nil, dirneuerr
	}

	content, conterr := os.ReadFile(filename)
	if conterr != nil {
		return nil, conterr
	}

	return content, nil
}
