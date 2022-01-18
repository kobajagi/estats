package stats

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const PathNet = "/proc/net/dev"

type NetProfile struct {
	ReceiveBytes  uint64
	TransmitBytes uint64
}

type NetUsage struct {
	f  *os.File
	in string
}

func NetUsageProvider(in string) *NetUsage {
	return &NetUsage{
		in: in,
	}
}

// Read calculates network speed for interface.
func (p *NetUsage) Read() ([]Stat, error) {
	if p.f == nil {
		f, err := os.Open(PathNet)
		if err != nil {
			return nil, err
		}

		p.f = f
	}

	content, err := readProcFile(p.f)
	if err != nil {
		return nil, err
	}

	net, err := p.parse(content)
	if err != nil {
		return nil, err
	}

	stats := []Stat{
		Stat{
			Name:   "download (bytes)",
			Metric: net.ReceiveBytes,
		},
		Stat{
			Name:   "upload (bytes)",
			Metric: net.TransmitBytes,
		},
	}

	return stats, nil
}

// parse parses network interface stats from /proc/net/dev.
func (p *NetUsage) parse(input []byte) (
	*NetProfile,
	error,
) {
	var net NetProfile
	buf := bytes.NewBuffer(input)
	scanner := bufio.NewScanner(buf)
	var n int
	found := false
	for scanner.Scan() {
		line := scanner.Text()
		//skip first two lines (header)
		n++
		if n <= 2 {
			continue
		}

		i, r, t, err := p.parseLine(line)
		if err != nil {
			return nil, err
		}

		if i == p.in {
			net.ReceiveBytes = r
			net.TransmitBytes = t
			found = true
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read /proc/net/dev: %w", err)
	}

	if !found {
		return nil, errors.New("network interface not found")
	}

	return &net, nil
}

// parseLine parses one network interface line from /proc/net/dev.
func (p *NetUsage) parseLine(line string) (
	i string,
	bReceive, bTransmit uint64,
	err error,
) {
	pattern := regexp.MustCompile(`\s+`)
	line = pattern.ReplaceAllString(line, " ")
	splits := strings.Split(strings.TrimSpace(line), " ")
	if len(splits) < 10 {
		err = errors.New("unsupported format of /proc/net/dev")
		return
	}

	i = strings.TrimSuffix(splits[0], ":")
	bReceive, err = strconv.ParseUint(splits[1], 10, 64)
	if err != nil {
		err = fmt.Errorf("NetUsage.parseLine: %w", err)
		return
	}
	bTransmit, err = strconv.ParseUint(splits[9], 10, 64)
	if err != nil {
		err = fmt.Errorf("NetUsage.parseLine: %w", err)
	}

	return
}
