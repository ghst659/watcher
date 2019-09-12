package watcher

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type tLog struct {
	count int
	t     *testing.T
}

func (tl *tLog) Print(v ...interface{}) {
	tl.t.Log(v...)
	tl.count++
}

func (tl *tLog) Printf(format string, v ...interface{}) {
	tl.t.Logf(format, v...)
	tl.count++
}

type verse struct {
	i     int
	lines []string
}

var stanza = &verse{
	i: 0,
	lines: []string{
		"In taberna quando sumus",
		"non curamus quid sit humus",
		"sed ad ludum properamus",
		"cui semper insudamus.",
		"Quid agatur in taberna",
		"ubi nummus est pincerna",
		"hoc est opus ut queratur",
		"sic quid loquar audiatur.",
	},
}

func (v *verse) reporter(ctx context.Context) ([]string, error) {
	lines := []string{
		fmt.Sprintf("%d: %q", v.i, v.lines[v.i%len(v.lines)]),
	}
	v.i++
	return lines, nil
}

type testTick struct {
	c chan time.Time
}

func newTicker() *testTick {
	return &testTick{
		c: make(chan time.Time),
	}
}

func (tick *testTick) Chan() <-chan time.Time {
	if tick == nil {
		return nil
	}
	return tick.c
}

func (tick *testTick) Stop() {
	if tick != nil {
		close(tick.c)
	}
}

func (tick *testTick) tick() {
	if tick != nil {
		tick.c <- time.Now()
	}
}

func TestWatch(t *testing.T) {
	// t.Log("START")
	nCtx, nCancel := context.WithCancel(context.Background())
	defer nCancel()
	clk := newTicker()
	out := &tLog{t: t}
	go Watch(nCtx, clk, stanza.reporter, out)
	for i := 0; i < 10; i++ {
		// t.Logf("SENDTICK: %d", i)
		clk.tick()
	}
	for out.count < 10 {
		// wait
	}
}
