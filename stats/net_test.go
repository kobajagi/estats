package stats

import (
	"os"
	"testing"
)

func TestNetParse(t *testing.T) {
	testdata, err := os.ReadFile("testdata/net_dev")
	if err != nil {
		t.Fatal(err)
	}

	p := NetUsageProvider("wlp2s0")
	net, err := p.parse(testdata)
	if err != nil {
		t.Error(err)
	}
	if net.ReceiveBytes != 311114747 {
		t.Errorf("expected net.ReceiveBytes=311114747, got %d", net.ReceiveBytes)
	}
	if net.TransmitBytes != 12117905 {
		t.Errorf("expected net.TransmitBytes=12117905, got %d", net.TransmitBytes)
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
