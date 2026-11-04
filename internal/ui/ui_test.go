package ui

import (
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/plyr4/tappy/internal/game"
)

func TestOnlyLeftButtonDownCounts(t *testing.T) {
	g := game.New(game.DefaultConfig, 0)
	m := New(g, nil)
	now := time.Now()
	m.now = func() time.Time {
		return now
	}
	send := func(msg tea.Msg) {
		r, _ := m.Update(msg)
		m = r.(Model)
	}

	left := tea.Mouse{Button: tea.MouseLeft}
	for _, msg := range []tea.Msg{
		tea.MouseReleaseMsg(left),
		tea.MouseMotionMsg(left),
		tea.MouseClickMsg(tea.Mouse{Button: tea.MouseRight}),
		tea.MouseWheelMsg(tea.Mouse{Button: tea.MouseWheelUp}),
		tea.KeyPressMsg{Code: 'a', Text: "a"},
	} {
		send(msg)
	}

	if g.Progress() != 0 {
		t.Fatalf("non-click events counted: %d", g.Progress())
	}

	for range 4 {
		send(tea.MouseClickMsg(left))
		send(tea.MouseReleaseMsg(left))
	}

	if g.Progress() != 4 {
		t.Fatalf("progress = %d, want 4", g.Progress())
	}
}

func TestSavesNewBest(t *testing.T) {
	g := game.New(game.DefaultConfig, 0)

	var saved []int

	save := func(b int) error {
		saved = append(saved, b)

		return nil
	}
	m := New(g, save)
	now := time.Now()
	m.now = func() time.Time {
		return now
	}
	left := tea.MouseClickMsg(tea.Mouse{Button: tea.MouseLeft})

	step := func(msg tea.Msg) {
		r, _ := m.Update(msg)
		m = r.(Model)
	}
	for range 5 {
		step(left)
	}

	now = now.Add(3 * time.Second)

	for range 7 {
		step(left)
	}

	now = now.Add(10 * time.Second)

	step(tickMsg{})
	step(tickMsg{})

	if len(saved) != 1 || saved[0] != 7 {
		t.Fatalf("saved = %v, want [7]", saved)
	}
}
