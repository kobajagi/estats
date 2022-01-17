package stats

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const PathStat = "/proc/stat"

type CpuProfile struct {
	User   uint64
	System uint64
	Idle   uint64
}

type CpuUsage struct {
	f         *os.File
	histTotal uint64
	histUsed  uint64
}

func CpuUsageProvider() *CpuUsage {
	return &CpuUsage{}
}

// Read calculates (total) CPU usage.
// Usage is calculated in a delta between current and previous
// polls.
func (p *CpuUsage) Read() ([]Stat, error) {
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

	cpu, err := p.parse(content)
	if err != nil {
		return nil, err
	}

	total := cpu.User + cpu.System + cpu.Idle
	used := cpu.User + cpu.System

	usage := int(math.Round(
		float64((used - p.histUsed) * 100 / (total - p.histTotal)),
	))
	p.histTotal = total
	p.histUsed = used

	return []Stat{
		{
			Name:   "cpu usage",
			Metric: fmt.Sprintf("%d%%", usage),
		},
	}, nil
}

// parse parses CPU stats from /proc/stat file.
// Field order is:
// user, nice, system, idle, iowait, irq, softirq, steal.
func (p *CpuUsage) parse(input []byte) (
	*CpuProfile,
	error,
) {
	var cpu CpuProfile
	buf := bytes.NewBuffer(input)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "cpu") {
			continue
		}

		name, user, system, idle, err := p.parseLine(line)
		if err != nil {
			return nil, err
		}
		if name == "cpu" {
			cpu.User = user
			cpu.System = system
			cpu.Idle = idle
			break
		}
	}

	return &cpu, nil
}

// parseLine parses one cpu line from /proc/stat.
func (p *CpuUsage) parseLine(line string) (
	cpu string,
	user, system, idle uint64,
	err error,
) {
	pattern := regexp.MustCompile(`\s+`)
	line = pattern.ReplaceAllString(line, " ")
	splits := strings.Split(line, " ")

	if len(splits) < 8 {
		err = errors.New("parser: malformed /proc/stat")
		return
	}

	cpu = splits[0]
	user, err = strconv.ParseUint(splits[1], 10, 64)
	if err != nil {
		return
	}
	system, err = strconv.ParseUint(splits[3], 10, 64)
	if err != nil {
		return
	}
	idle, err = strconv.ParseUint(splits[4], 10, 64)
	if err != nil {
		return
	}

	return
}
