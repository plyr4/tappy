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
	return filepath.Join(dir, "tapper", "best.json"), nil
}

// Load returns the stored personal best, or 0 if missing or malformed.
func Load(path string) int {
	b, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	var d data
	if json.Unmarshal(b, &d) != nil || d.Best < 0 {
		return 0
	}
	return d.Best
}

// Save atomically writes the personal best to path.
func Save(path string, best int) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.Marshal(data{Best: best})
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
