package stats

import (
	"fmt"
	"math"
	"syscall"
)

type DiskProfile struct {
	Size, Available uint64
}

type DiskUsage struct {
	path string
}

func DiskUsageProvider(path string) *DiskUsage {
	return &DiskUsage{path: path}
}

// Read calculates disk usage.
func (p *DiskUsage) Read() ([]Stat, error) {
	disk, err := p.syscall()
	if err != nil {
		return nil, err
	}

	usage := math.Round(float64((disk.Size - disk.Available) * 100 / disk.Size))

	return []Stat{
		{
			Name:   "disk usage (%)",
			Metric: int(usage),
		},
	}, nil
}

// sysctll calls statfs to get disk information
func (p *DiskUsage) syscall() (*DiskProfile, error) {
	var disk DiskProfile
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(p.path, &fs)
	if err != nil {
		return nil, fmt.Errorf("statfs syscall failed: %w", err)
	}

	disk.Size = fs.Blocks * uint64(fs.Bsize)
	disk.Available = fs.Bfree * uint64(fs.Bsize)

	return &disk, nil
}
