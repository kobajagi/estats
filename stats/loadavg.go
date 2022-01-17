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

// Read fetches load avarage statistics.
func (p *LoadAvg) Read() (map[string]string, error) {
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
	map[string]string,
	error,
) {
	splits := strings.Split(string(input), " ")
	// We will only take first 3, others may be os/kernel specific.
	if len(splits) < 3 {
		return nil, errors.New("parser: malformed /proc/loadavg")
	}

	output := map[string]string{}
	keys := []string{"1m", "5m", "15m"}

	for index, field := range splits[:3] {
		output["loadavg."+keys[index]] = field
	}

	return output, nil
}
