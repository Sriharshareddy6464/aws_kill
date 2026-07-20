package utils

import (
	"encoding/json"
	"io"
	"os"
)

// ReadJSON loads a struct from a JSON file path
func ReadJSON(path string, target interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, target)
}

// WriteJSON saves a struct to a JSON file path
func WriteJSON(path string, data interface{}) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, bytes, 0644)
}
