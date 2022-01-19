package stats

import (
	"os"
	"testing"
)

func TestMemParse(t *testing.T) {
	testdata, err := os.ReadFile("testdata/meminfo")
	if err != nil {
		t.Fatal(err)
	}

	p := MemUsageProvider()
	mem, err := p.parse(testdata)
	if err != nil {
		t.Error(err)
	}
	if mem.Total != 3979204 {
		t.Errorf("expected mem.Total=113064, got %d", mem.Total)
	}
	if mem.Free != 1461588 {
		t.Errorf("expected mem.Free=1461588, got %d", mem.Free)
	}
	if mem.Available != 2647256 {
		t.Errorf("expected mem.Available=2647256, got %d", mem.Available)
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
