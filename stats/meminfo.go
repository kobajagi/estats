package stats

import (
	"bufio"
	"bytes"
	"errors"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const PathMeminfo = "/proc/meminfo"

type MemProfile struct {
	Total, Free, Available uint64
}

type MemUsage struct {
	f *os.File
}

func MemUsageProvider() *MemUsage {
	return &MemUsage{}
}

// Read calculates memory usage.
func (p *MemUsage) Read() ([]Stat, error) {
	if p.f == nil {
		f, err := os.Open(PathMeminfo)
		if err != nil {
			return nil, err
		}

		p.f = f
	}

	content, err := readProcFile(p.f)
	if err != nil {
		return nil, err
	}

	mem, err := p.parse(content)
	if err != nil {
		return nil, err
	}

	usage := math.Round(float64(mem.Total-mem.Available) * 100 / float64(mem.Total))

	return []Stat{
		Stat{
			Name:   "memory usage (%)",
			Metric: int(usage),
		},
	}, nil
}

// parse parses memory stats from /proc/meminfo.
func (p *MemUsage) parse(input []byte) (
	*MemProfile,
	error,
) {
	mem := MemProfile{}
	buf := bytes.NewBuffer(input)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := scanner.Text()
		key, value, err := p.parseLine(line)
		if err != nil {
			return nil, err
		}

		switch key {
		case "MemTotal":
			mem.Total = value
		case "MemFree":
			mem.Free = value
		case "MemAvailable":
			mem.Available = value
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &mem, nil
}

// parseLine parses single line from /proc/meminfo.
func (p *MemUsage) parseLine(line string) (
	key string,
	value uint64,
	err error,
) {
	pattern := regexp.MustCompile(`\s+`)
	line = pattern.ReplaceAllString(line, " ")
	splits := strings.Split(line, " ")
	if len(splits) < 2 {
		err = errors.New("usupported format of /proc/meminfo")
		return
	}

	key = strings.TrimSuffix(splits[0], ":")
	value, err = strconv.ParseUint(splits[1], 10, 64)

	return
}
