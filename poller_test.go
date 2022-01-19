package main

import (
	"testing"

	"github.com/kobajagi/estats/output"
	"github.com/kobajagi/estats/stats"
)

var TestStats = []stats.Stat{
	{
		Name:   "test1",
		Metric: "val1",
	},
	{
		Name:   "test2",
		Metric: "val2",
	},
	{
		Name:   "test3",
		Metric: "val3",
	},
}

type TestWriter struct {
	r chan *output.Report
}

func (w *TestWriter) Write(r output.Report) error {
	w.r <- &r
	return nil
}

type TestProvider struct {
	r []stats.Stat
}

func (p TestProvider) Read() ([]stats.Stat, error) {
	return p.r, nil
}

func TestPoll(t *testing.T) {
	p := NewPoller(
		TestProvider{TestStats[0:2]},
		TestProvider{[]stats.Stat{TestStats[2]}},
	)
	c := make(chan *output.Report)
	w := TestWriter{r: c}
	p.writer = &w

	p.Poll()
	r := <-c

	if len(r.Metrics) != 3 {
		t.Errorf("expected 3 elements in Report, got %d", len(r.Metrics))
	}
}
