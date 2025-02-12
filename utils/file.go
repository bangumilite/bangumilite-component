package utils

import (
	"encoding/json"
	"os"
)

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
