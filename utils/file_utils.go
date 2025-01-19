package utils

import (
	"encoding/json"
	"os"
)

// SaveToJSON saves any struct type to a local JSON file.
func SaveToJSON[T any](data T, fn string) error {
	file, err := os.Create(fn)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	return nil
}
