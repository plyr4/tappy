package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "best.json")

	got := Load(path)
	if got != 0 {
		t.Fatalf("missing file: got %d", got)
	}

	err := Save(path, 92)
	if err != nil {
		t.Fatal(err)
	}

	got = Load(path)
	if got != 92 {
		t.Fatalf("got %d, want 92", got)
	}
}

func TestLoadMalformed(t *testing.T) {
	path := filepath.Join(t.TempDir(), "best.json")

	for _, s := range []string{"not json", `{"best":-3}`, ""} {
		err := os.WriteFile(path, []byte(s), 0o600)
		if err != nil {
			t.Fatal(err)
		}

		got := Load(path)
		if got != 0 {
			t.Fatalf("%q: got %d", s, got)
		}
	}
}
