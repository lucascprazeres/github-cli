package file

import (
	"encoding/json"
	"os"
)

func Write(filename string, content any) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(content)
}
