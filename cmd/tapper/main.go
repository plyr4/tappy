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
		fmt.Fprintln(os.Stderr, "tapper:", err)
		os.Exit(1)
	}
	g := game.New(game.DefaultConfig, storage.Load(path))
	m := ui.New(g, func(best int) error { return storage.Save(path, best) })
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "tapper:", err)
		os.Exit(1)
	}
}
