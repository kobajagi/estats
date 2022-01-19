package stats

import (
	"os"
	"testing"
)

func TestLoadavgParse(t *testing.T) {
	testdata, err := os.ReadFile("testdata/loadavg")
	if err != nil {
		t.Fatal(err)
	}

	p := LoadAvgProvider()
	lavg, err := p.parse(testdata)
	if err != nil {
		t.Error(err)
	}
	if lavg["1m"] != "0.27" {
		t.Errorf("expected loadavg.1m=0.27, got %s", lavg["1m"])
	}
	if lavg["5m"] != "0.42" {
		t.Errorf("expected loadavg.5m=0.42, got %s", lavg["5m"])
	}
	if lavg["15m"] != "0.64" {
		t.Errorf("expected loadavg.15m=0.64, got %s", lavg["15m"])
	}

	baddata, err := os.ReadFile("testdata/stat")
	if err != nil {
		t.Fatal(err)
	}
	_, err = p.parse(baddata)
	if err == nil {
		t.Error("expected error on bad data, got nil")
	}
}
