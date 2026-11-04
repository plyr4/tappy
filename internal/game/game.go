// Package game implements the tappy state machine.
package game

import "time"

// Phase is a game state.
type Phase int

// Game phases.
const (
	Armed Phase = iota
	Countdown
	Playing
	Results
)

// Config holds game timing and trigger settings.
type Config struct {
	TriggerClicks int
	TriggerWindow time.Duration
	CountdownFrom int
	GameDuration  time.Duration
	ResultsDelay  time.Duration
}

// DefaultConfig is the standard game configuration.
var DefaultConfig = Config{
	TriggerClicks: 5,
	TriggerWindow: time.Second,
	CountdownFrom: 3,
	GameDuration:  10 * time.Second,
	ResultsDelay:  3 * time.Second,
}

// Result is the outcome of a round.
type Result struct {
	Clicks   int
	Duration time.Duration
	Best     int
	NewBest  bool
}

// CPS returns clicks per second over the measured duration.
func (r Result) CPS() float64 {
	if r.Duration <= 0 {
		return 0
	}

	clicks := float64(r.Clicks)
	seconds := r.Duration.Seconds()

	return clicks / seconds
}

// Game tracks the current phase, clicks, and personal best.
type Game struct {
	cfg    Config
	phase  Phase
	taps   []time.Time
	start  time.Time
	clicks int
	best   int
	result *Result
}

// New returns a game in armed mode.
func New(cfg Config, best int) *Game {
	return &Game{cfg: cfg, best: best}
}

// Config returns the game configuration.
func (g *Game) Config() Config {
	return g.cfg
}

// Phase returns the current phase.
func (g *Game) Phase() Phase {
	return g.phase
}

// Clicks returns the clicks counted in the current round.
func (g *Game) Clicks() int {
	return g.clicks
}

// Best returns the personal best.
func (g *Game) Best() int {
	return g.best
}

// Progress returns the clicks counted toward the trigger.
func (g *Game) Progress() int {
	return len(g.taps)
}

// Result returns the last round result, or nil.
func (g *Game) Result() *Result {
	return g.result
}

// Click registers a left button-down at now.
func (g *Game) Click(now time.Time) {
	g.Advance(now)

	switch g.phase {
	case Armed:
		g.taps = append(g.taps, now)
		if len(g.taps) >= g.cfg.TriggerClicks {
			g.taps = nil
			g.phase = Countdown
			g.start = now
		}
	case Playing:
		g.clicks++
	}
}

// Advance moves the game through any time-based transitions up to now.
func (g *Game) Advance(now time.Time) {
	for {
		switch g.phase {
		case Armed:
			g.prune(now)

			return
		case Countdown:
			goAt := g.start.Add(time.Duration(g.cfg.CountdownFrom) * time.Second)
			if now.Before(goAt) {
				return
			}

			g.phase = Playing
			g.start = goAt
			g.clicks = 0
		case Playing:
			end := g.start.Add(g.cfg.GameDuration)
			if now.Before(end) {
				return
			}

			g.finish(end)
		case Results:
			if now.Sub(g.start) < g.cfg.ResultsDelay {
				return
			}

			g.phase = Armed
			g.taps = nil
		}
	}
}

// Abort cancels any round in progress and rearms.
func (g *Game) Abort() {
	if g.phase == Countdown || g.phase == Playing {
		g.phase = Armed
		g.clicks = 0
	}

	g.taps = nil
}

// CountdownValue returns the countdown number to display.
func (g *Game) CountdownValue(now time.Time) int {
	elapsed := int(now.Sub(g.start) / time.Second)
	n := g.cfg.CountdownFrom - elapsed
	n = min(n, g.cfg.CountdownFrom)

	return max(n, 1)
}

// Remaining returns the time left in the current round.
func (g *Game) Remaining(now time.Time) time.Duration {
	end := g.start.Add(g.cfg.GameDuration)
	remaining := end.Sub(now)

	return max(remaining, 0)
}

func (g *Game) finish(end time.Time) {
	r := Result{Clicks: g.clicks, Duration: end.Sub(g.start)}
	if r.Clicks > g.best {
		g.best = r.Clicks
		r.NewBest = true
	}

	r.Best = g.best
	g.result = &r
	g.phase = Results
	g.start = end
}

func (g *Game) prune(now time.Time) {
	i := 0
	for i < len(g.taps) && now.Sub(g.taps[i]) >= g.cfg.TriggerWindow {
		i++
	}

	g.taps = g.taps[i:]
}
