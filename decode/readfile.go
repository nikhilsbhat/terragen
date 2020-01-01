package decode

import (
	"io/ioutil"
	"os"
)

// ReadFile as name specifies it just reads the content of file ans returns it.
func ReadFile(filename string) ([]byte, error) {
	if _, dirneuerr := os.Stat(filename); os.IsNotExist(dirneuerr) {
		return nil, dirneuerr
	}

	content, conterr := ioutil.ReadFile(filename)
	if conterr != nil {
		return nil, conterr
	}
	return content, nil
	//return nil, nil
}
