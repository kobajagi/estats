package stats

import (
	"errors"
	"os"
	"strings"
)

const PathLoadAvg = "/proc/loadavg"

type LoadAvg struct {
	f *os.File
}

func LoadAvgProvider() *LoadAvg {
	return &LoadAvg{}
}

// Read reads load average.
func (p *LoadAvg) Read() ([]Stat, error) {
	if p.f == nil {
		f, err := os.Open(PathLoadAvg)
		if err != nil {
			return nil, err
		}

		p.f = f
	}

	content, err := readProcFile(p.f)
	if err != nil {
		return nil, err
	}

	stats, err := p.parse(content)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

// parse parses content of /proc/loadavg file
func (p *LoadAvg) parse(input []byte) (
	[]Stat,
	error,
) {
	splits := strings.Split(string(input), " ")
	if len(splits) < 3 {
		return nil, errors.New("parser: malformed /proc/loadavg")
	}

	stats := make([]Stat, 0, 3)
	keys := []string{"1m", "5m", "15m"}

	for index, field := range splits[:3] {
		stats = append(stats, Stat{
			Name:   "loadavg " + keys[index],
			Metric: field,
		})
	}

	return stats, nil
}
