package output

import "github.com/kobajagi/estats/stats"

type OutputWriter interface {
	Write(Report) error
}
type Report struct {
	Timestamp int64        `json:"timestamp"`
	Metrics   []stats.Stat `json:"metrics"`
}
