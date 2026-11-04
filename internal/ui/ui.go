// Package ui renders the game with Bubble Tea.
package ui

import (
	"fmt"
	"strconv"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/plyr4/tappy/internal/game"
)

const tickRate = 50 * time.Millisecond

var (
	big   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	faint = lipgloss.NewStyle().Faint(true)
	good  = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
)

type tickMsg struct{}

// Model is the Bubble Tea model.
type Model struct {
	game    *game.Game
	save    func(int) error
	saved   int
	saveErr error
	now     func() time.Time
	w, h    int
}

// New returns a Bubble Tea model for g that persists new bests with save.
func New(g *game.Game, save func(int) error) Model {
	return Model{game: g, save: save, saved: g.Best(), now: time.Now}
}

// Init starts the tick loop.
func (m Model) Init() tea.Cmd {
	return tick()
}

// tick schedules the next tickMsg.
func tick() tea.Cmd {
	return tea.Tick(tickRate, newTickMsg)
}

// newTickMsg returns a tickMsg.
func newTickMsg(time.Time) tea.Msg {
	return tickMsg{}
}

// Update handles input and ticks.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return m, tea.Quit
		case "ctrl+z":
			m.game.Abort()

			return m, tea.Suspend
		}
	case tea.MouseClickMsg:
		if msg.Button == tea.MouseLeft {
			m.game.Click(m.now())
		}
	case tea.BlurMsg, tea.ResumeMsg:
		m.game.Abort()
	case tickMsg:
		m.game.Advance(m.now())

		cmd = tick()
	}

	best := m.game.Best()
	if best > m.saved && m.save != nil {
		m.saveErr = m.save(best)
		m.saved = best
	}

	return m, cmd
}

// View renders the current game state.
func (m Model) View() tea.View {
	content := lipgloss.Place(m.w, m.h, lipgloss.Center, lipgloss.Center, m.content())

	v := tea.NewView(content)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	v.ReportFocus = true

	return v
}

// content returns the text for the current phase.
func (m Model) content() string {
	var lines []string

	switch m.game.Phase() {
	case game.Armed:
		lines = m.armed()
	case game.Countdown:
		lines = m.countdown()
	case game.Playing:
		lines = m.playing()
	case game.Results:
		lines = results(m.game.Result())
	}

	if m.saveErr != nil {
		msg := faint.Render("could not save best: " + m.saveErr.Error())
		lines = append(lines, "", msg)
	}

	return lipgloss.JoinVertical(lipgloss.Center, lines...)
}

// armed returns the lines for armed mode.
func (m Model) armed() []string {
	n := m.game.Config().TriggerClicks

	var lines []string

	r := m.game.Result()
	if r != nil {
		prompt := fmt.Sprintf("tap %d times to play again", n)
		lines = append(results(r), "", faint.Render(prompt))
	} else {
		prompt := fmt.Sprintf("tap %d times to start", n)
		lines = []string{prompt}
	}

	progress := fmt.Sprintf("%d / %d", m.game.Progress(), n)

	return append(lines, faint.Render(progress))
}

// countdown returns the lines for the countdown.
func (m Model) countdown() []string {
	value := m.game.CountdownValue(m.now())

	return []string{big.Render(strconv.Itoa(value))}
}

// playing returns the lines for an active round.
func (m Model) playing() []string {
	count := strconv.Itoa(m.game.Clicks())
	if m.game.Clicks() == 0 {
		count = "GO!"
	}

	remaining := m.game.Remaining(m.now()).Seconds()
	timer := fmt.Sprintf("%.2fs", remaining)

	return []string{big.Render(count), "", faint.Render(timer)}
}

// results returns the lines describing r.
func results(r *game.Result) []string {
	clicks := fmt.Sprintf("%d clicks", r.Clicks)
	cps := fmt.Sprintf("%.1f clicks/sec", r.CPS())
	best := fmt.Sprintf("best: %d", r.Best)

	lines := []string{big.Render(clicks), cps, faint.Render(best)}

	if r.NewBest {
		lines = append(lines, good.Render("new best!"))
	}

	return lines
}
