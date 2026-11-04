package ui

import (
	"fmt"
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

func (m Model) Init() tea.Cmd { return tick() }

func tick() tea.Cmd {
	return tea.Tick(tickRate, func(time.Time) tea.Msg { return tickMsg{} })
}

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
	if b := m.game.Best(); b > m.saved && m.save != nil {
		m.saveErr = m.save(b)
		m.saved = b
	}
	return m, cmd
}

func (m Model) View() tea.View {
	v := tea.NewView(lipgloss.Place(m.w, m.h, lipgloss.Center, lipgloss.Center, m.content()))
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	v.ReportFocus = true
	return v
}

func (m Model) content() string {
	g := m.game
	now := m.now()
	n := g.Config().TriggerClicks
	var lines []string
	switch g.Phase() {
	case game.Armed:
		if r := g.Result(); r != nil {
			lines = append(results(r), "", faint.Render(fmt.Sprintf("tap %d times to play again", n)))
		} else {
			lines = []string{fmt.Sprintf("tap %d times to start", n)}
		}
		lines = append(lines, faint.Render(fmt.Sprintf("%d / %d", g.Progress(), n)))
	case game.Countdown:
		lines = []string{big.Render(fmt.Sprint(g.CountdownValue(now)))}
	case game.Playing:
		count := fmt.Sprint(g.Clicks())
		if g.Clicks() == 0 {
			count = "GO!"
		}
		lines = []string{big.Render(count), "", faint.Render(fmt.Sprintf("%.2fs", g.Remaining(now).Seconds()))}
	case game.Results:
		lines = results(g.Result())
	}
	if m.saveErr != nil {
		lines = append(lines, "", faint.Render("could not save best: "+m.saveErr.Error()))
	}
	return lipgloss.JoinVertical(lipgloss.Center, lines...)
}

func results(r *game.Result) []string {
	lines := []string{
		big.Render(fmt.Sprintf("%d clicks", r.Clicks)),
		fmt.Sprintf("%.1f clicks/sec", r.CPS()),
		faint.Render(fmt.Sprintf("best: %d", r.Best)),
	}
	if r.NewBest {
		lines = append(lines, good.Render("new best!"))
	}
	return lines
}
