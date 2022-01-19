package stats

import (
	"errors"
	"os"
	"regexp"
	"strings"
)

var RegexLoadAvg = regexp.MustCompile(`[0-9]*\.[0-9]{2}`)

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

	lavg, err := p.parse(content)
	if err != nil {
		return nil, err
	}

	stats := make([]Stat, 0, len(lavg))
	for k, v := range lavg {
		stats = append(stats, Stat{
			Name:   "loadavg " + k,
			Metric: v,
		})
	}

	return stats, nil
}

// parse parses content of /proc/loadavg file
func (p *LoadAvg) parse(input []byte) (
	map[string]string,
	error,
) {
	splits := strings.Split(string(input), " ")
	if len(splits) < 3 {
		return nil, errors.New("parser: malformed /proc/loadavg")
	}

	lavg := map[string]string{}
	keys := []string{"1m", "5m", "15m"}

	for index, field := range splits[:3] {
		if !RegexLoadAvg.MatchString(field) {
			return nil, errors.New("parser: malformed /proc/loadavg")
		}
		lavg[keys[index]] = field
	}

	return lavg, nil
}
