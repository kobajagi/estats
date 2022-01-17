package stats

import (
	"io"
	"os"
)

// Provider collects a single type of statistic.
type Provider interface {
	Read() ([]Stat, error)
}

// Stat is a single metric.
type Stat struct {
	Name   string `json:"name"`
	Metric string `json:"metric"`
}

// reads proc file content, reseting pointer aferwards.
func readProcFile(f *os.File) ([]byte, error) {
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	return content, err
}
