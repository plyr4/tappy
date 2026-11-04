package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "best.json")
	if got := Load(path); got != 0 {
		t.Fatalf("missing file: got %d", got)
	}
	if err := Save(path, 92); err != nil {
		t.Fatal(err)
	}
	if got := Load(path); got != 92 {
		t.Fatalf("got %d, want 92", got)
	}
}

func TestLoadMalformed(t *testing.T) {
	path := filepath.Join(t.TempDir(), "best.json")
	for _, s := range []string{"not json", `{"best":-3}`, ""} {
		os.WriteFile(path, []byte(s), 0o644)
		if got := Load(path); got != 0 {
			t.Fatalf("%q: got %d", s, got)
		}
	}
}
