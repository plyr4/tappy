package game

import "time"

type Phase int

const (
	Armed Phase = iota
	Countdown
	Playing
	Results
)

type Config struct {
	TriggerClicks int
	TriggerWindow time.Duration
	CountdownFrom int
	GameDuration  time.Duration
	ResultsDelay  time.Duration
}

var DefaultConfig = Config{
	TriggerClicks: 5,
	TriggerWindow: time.Second,
	CountdownFrom: 3,
	GameDuration:  10 * time.Second,
	ResultsDelay:  3 * time.Second,
}

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
	return float64(r.Clicks) / r.Duration.Seconds()
}

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

func (g *Game) Config() Config  { return g.cfg }
func (g *Game) Phase() Phase    { return g.phase }
func (g *Game) Clicks() int     { return g.clicks }
func (g *Game) Best() int       { return g.best }
func (g *Game) Progress() int   { return len(g.taps) }
func (g *Game) Result() *Result { return g.result }

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
	n := g.cfg.CountdownFrom - int(now.Sub(g.start)/time.Second)
	return max(1, min(n, g.cfg.CountdownFrom))
}

// Remaining returns the time left in the current round.
func (g *Game) Remaining(now time.Time) time.Duration {
	return max(0, g.start.Add(g.cfg.GameDuration).Sub(now))
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
