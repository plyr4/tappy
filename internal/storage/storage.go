// Package storage persists the personal best.
package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type data struct {
	Best int `json:"best"`
}

// Path returns the personal best file in the user's application data directory.
func Path() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(dir, "tappy", "best.json")

	return path, nil
}

// Load returns the stored personal best, or 0 if missing or malformed.
func Load(path string) int {
	b, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return 0
	}

	var d data

	err = json.Unmarshal(b, &d)
	if err != nil {
		return 0
	}

	if d.Best < 0 {
		return 0
	}

	return d.Best
}

// Save atomically writes the personal best to path.
func Save(path string, best int) error {
	err := os.MkdirAll(filepath.Dir(path), 0o750)
	if err != nil {
		return err
	}

	b, err := json.Marshal(data{Best: best})
	if err != nil {
		return err
	}

	tmp := path + ".tmp"

	err = os.WriteFile(tmp, b, 0o600)
	if err != nil {
		return err
	}

	return os.Rename(tmp, path)
}
