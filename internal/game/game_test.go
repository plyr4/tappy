package game

import (
	"testing"
	"time"
)

var t0 = time.Now()

func at(d time.Duration) time.Time {
	return t0.Add(d)
}

func ms(n int) time.Duration {
	return time.Duration(n) * time.Millisecond
}

func tapN(g *Game, start time.Duration, n int, gap time.Duration) {
	for i := range n {
		g.Click(at(start + time.Duration(i)*gap))
	}
}

func started(t *testing.T) *Game {
	t.Helper()

	g := New(DefaultConfig, 0)
	tapN(g, 0, 5, ms(100))

	if g.Phase() != Countdown {
		t.Fatalf("phase = %v, want Countdown", g.Phase())
	}

	return g
}

func TestTrigger(t *testing.T) {
	g := New(DefaultConfig, 0)
	for i := range 4 {
		g.Click(at(ms(i * 100)))

		if g.Progress() != i+1 || g.Phase() != Armed {
			t.Fatalf("after %d clicks: progress=%d phase=%v", i+1, g.Progress(), g.Phase())
		}
	}

	g.Click(at(ms(400)))

	if g.Phase() != Countdown || g.Progress() != 0 {
		t.Fatalf("phase=%v progress=%d, want Countdown 0", g.Phase(), g.Progress())
	}
}

func TestTriggerWindowExpires(t *testing.T) {
	g := New(DefaultConfig, 0)
	tapN(g, 0, 4, ms(100))
	g.Advance(at(ms(1350)))

	if g.Progress() != 0 {
		t.Fatalf("progress = %d, want 0 after window expires", g.Progress())
	}

	g.Click(at(ms(1400)))

	if g.Phase() != Armed || g.Progress() != 1 {
		t.Fatalf("phase=%v progress=%d, want Armed 1", g.Phase(), g.Progress())
	}
}

func TestTriggerRollingWindow(t *testing.T) {
	g := New(DefaultConfig, 0)
	tapN(g, 0, 5, ms(300))

	if g.Phase() != Armed {
		t.Fatalf("5 clicks over 1.2s triggered")
	}

	if g.Progress() != 4 {
		t.Fatalf("progress = %d, want 4", g.Progress())
	}

	g.Click(at(ms(1250)))

	if g.Phase() != Countdown {
		t.Fatalf("5 clicks within 1s did not trigger")
	}
}

func TestCountdown(t *testing.T) {
	g := started(t)

	trig := ms(400)
	for _, c := range []struct {
		d    time.Duration
		want int
	}{{0, 3}, {ms(999), 3}, {ms(1000), 2}, {ms(2500), 1}} {
		g.Advance(at(trig + c.d))

		if g.Phase() != Countdown || g.CountdownValue(at(trig+c.d)) != c.want {
			t.Fatalf("at %v: phase=%v value=%d, want %d", c.d, g.Phase(), g.CountdownValue(at(trig+c.d)), c.want)
		}
	}

	g.Click(at(trig + ms(2999)))
	g.Advance(at(trig + ms(3000)))

	if g.Phase() != Playing || g.Clicks() != 0 {
		t.Fatalf("phase=%v clicks=%d, want Playing 0", g.Phase(), g.Clicks())
	}

	if g.Remaining(at(trig+ms(3000))) != 10*time.Second {
		t.Fatalf("remaining = %v", g.Remaining(at(trig+ms(3000))))
	}
}

func TestGameCountsAndCompletes(t *testing.T) {
	g := started(t)
	goAt := ms(3400)
	g.Click(at(goAt))
	tapN(g, goAt+ms(10), 86, ms(100))
	g.Click(at(goAt + 10*time.Second - time.Nanosecond))

	if g.Phase() != Playing || g.Clicks() != 88 {
		t.Fatalf("phase=%v clicks=%d, want Playing 88", g.Phase(), g.Clicks())
	}

	g.Click(at(goAt + 10*time.Second))

	if g.Phase() != Results {
		t.Fatalf("phase = %v, want Results", g.Phase())
	}

	r := g.Result()
	if r.Clicks != 88 || r.Duration != 10*time.Second {
		t.Fatalf("result = %+v", r)
	}

	if r.CPS() != 8.8 {
		t.Fatalf("cps = %v, want 8.8", r.CPS())
	}

	if !r.NewBest || r.Best != 88 || g.Best() != 88 {
		t.Fatalf("best not updated: %+v", r)
	}
}

func TestLateTickDoesNotExtendGame(t *testing.T) {
	g := started(t)
	g.Advance(at(ms(3400) + 30*time.Second))

	if g.Result() == nil || g.Result().Duration != 10*time.Second {
		t.Fatalf("result = %+v", g.Result())
	}
}

func TestResultsRearm(t *testing.T) {
	g := New(DefaultConfig, 50)
	tapN(g, 0, 5, ms(100))
	end := ms(3400) + 10*time.Second
	g.Advance(at(end))
	g.Click(at(end + ms(500)))

	if g.Phase() != Results || g.Progress() != 0 {
		t.Fatalf("clicks during results counted: phase=%v progress=%d", g.Phase(), g.Progress())
	}

	r := g.Result()
	if r.NewBest || r.Best != 50 {
		t.Fatalf("result = %+v, want best 50 not new", r)
	}

	g.Advance(at(end + DefaultConfig.ResultsDelay))

	if g.Phase() != Armed || g.Best() != 50 {
		t.Fatalf("phase=%v best=%d", g.Phase(), g.Best())
	}

	tapN(g, end+5*time.Second, 5, ms(50))

	if g.Phase() != Countdown {
		t.Fatalf("did not rearm")
	}
}

func TestAbort(t *testing.T) {
	g := started(t)
	g.Advance(at(ms(4000)))
	g.Click(at(ms(4100)))
	g.Abort()

	if g.Phase() != Armed || g.Clicks() != 0 || g.Result() != nil {
		t.Fatalf("abort: phase=%v clicks=%d", g.Phase(), g.Clicks())
	}
}

func TestCPSZeroDuration(t *testing.T) {
	if (Result{Clicks: 5}).CPS() != 0 {
		t.Fatal("expected 0")
	}
}
