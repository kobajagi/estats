package stats

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const PathStat = "/proc/stat"

type CpuUsage struct {
	f         *os.File
	histTotal uint64
	histUsed  uint64
}

func CpuUsageProvider() *CpuUsage {
	return &CpuUsage{}
}

// Read derives CPU usage.
func (p *CpuUsage) Read() (map[string]string, error) {
	if p.f == nil {
		f, err := os.Open(PathStat)
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

// parse parses content of /proc/stat file.
// Only row with total CPU statistics is used.
// Field order is:
// user, nice, system, idle, iowait, irq, softirq, steal.
// Usage percent calculated as:
// (user+system)*100/user+system+idle
func (p *CpuUsage) parse(input []byte) (
	map[string]string,
	error,
) {
	var row string
	buf := bytes.NewBuffer(input)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu ") {
			row = strings.Trim(
				strings.TrimPrefix(line, "cpu"),
				" \t",
			)
			break
		}
	}

	if row == "" {
		return nil, errors.New("parser: malformed /proc/stat")
	}

	splits := strings.Split(string(row), " ")
	if len(splits) < 7 {
		return nil, errors.New("parser: malformed /proc/stat")
	}

	vals := []uint64{}
	for _, field := range splits[:7] {
		val, err := strconv.ParseUint(field, 10, 64)
		if err != nil {
			return nil, errors.New("parser: malformed /proc/stat")
		}

		vals = append(vals, val)
	}

	total := vals[0] + vals[2] + vals[3]
	used := vals[0] + vals[2]

	usage := int((used - p.histUsed) * 100 / (total - p.histTotal))
	p.histTotal = total
	p.histUsed = used

	return map[string]string{
		"cpu.percent": fmt.Sprintf("%d", usage),
	}, nil
}
