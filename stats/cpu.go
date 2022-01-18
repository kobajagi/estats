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
// Usage is calculated using delta between current and previous
// read.
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
	usage := math.Round(float64(used-p.histUsed) * 100 / float64(total-p.histTotal))
	p.histTotal = total
	p.histUsed = used

	return []Stat{
		{
			Name:   "cpu usage (%)",
			Metric: int(usage),
		},
	}, nil
}

// parse parses CPU stats from /proc/stat file.
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
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &cpu, nil
}

// parseLine parses one cpu line from /proc/stat.
// Field order is:
// cpu name, user, nice, system, idle, iowait, irq, softirq, steal.
func (p *CpuUsage) parseLine(line string) (
	cpu string,
	user, system, idle uint64,
	err error,
) {
	pattern := regexp.MustCompile(`\s+`)
	line = pattern.ReplaceAllString(line, " ")
	splits := strings.Split(line, " ")
	if len(splits) < 8 {
		err = errors.New("unsupported format of /proc/stat")
		return
	}

	cpu = splits[0]
	user, err = strconv.ParseUint(splits[1], 10, 64)
	if err != nil {
		err = fmt.Errorf("CpuUsage.parseLine: %w", err)
		return
	}
	system, err = strconv.ParseUint(splits[3], 10, 64)
	if err != nil {
		err = fmt.Errorf("CpuUsage.parseLine: %w", err)
		return
	}
	idle, err = strconv.ParseUint(splits[4], 10, 64)
	if err != nil {
		err = fmt.Errorf("CpuUsage.parseLine: %w", err)
	}

	return
}
