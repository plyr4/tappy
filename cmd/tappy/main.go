// Command tappy is a terminal clicking-speed game.
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/plyr4/tappy/internal/game"
	"github.com/plyr4/tappy/internal/storage"
	"github.com/plyr4/tappy/internal/ui"
)

func main() {
	path, err := storage.Path()
	if err != nil {
		exit(err)
	}

	best := storage.Load(path)
	g := game.New(game.DefaultConfig, best)

	save := func(best int) error {
		return storage.Save(path, best)
	}

	m := ui.New(g, save)

	_, err = tea.NewProgram(m).Run()
	if err != nil {
		exit(err)
	}
}

// exit prints err and exits with status 1.
func exit(err error) {
	fmt.Fprintln(os.Stderr, "tappy:", err)
	os.Exit(1)
}
