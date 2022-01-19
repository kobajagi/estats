package stats

import (
	"os"
	"testing"
)

func TestCpuParse(t *testing.T) {
	testdata, err := os.ReadFile("testdata/stat")
	if err != nil {
		t.Fatal(err)
	}

	p := CpuUsageProvider()
	cpu, err := p.parse(testdata)
	if err != nil {
		t.Error(err)
	}
	if cpu.User != 113064 {
		t.Errorf("expected cpu.User=113064, got %d", cpu.User)
	}
	if cpu.System != 39324 {
		t.Errorf("expected cpu.System=39324, got %d", cpu.System)
	}
	if cpu.Idle != 1882090 {
		t.Errorf("expected cpu.Idle=1882090, got %d", cpu.Idle)
	}

	baddata, err := os.ReadFile("testdata/loadavg")
	if err != nil {
		t.Fatal(err)
	}
	_, err = p.parse(baddata)
	if err == nil {
		t.Error("expected error on bad data, got nil")
	}
}
