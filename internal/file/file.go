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

func Read(filename string) (map[string]interface{}, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var content map[string]interface{}
	err = json.Unmarshal(data, &content)
	if err != nil {
		return nil, err
	}

	return content, err
}
