package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kobajagi/estats/stats"
)

const HelpString = `estats is a system metrics collection system which polls:
- Load average values
- Derived CPU percentage values
- Network interface statistics (in -n option set)
- Disk partition usage in percent (if -p option set)
- Memory usage in percent`

func main() {
	flag.Usage = func() {
		if len(os.Args) == 2 && os.Args[1] == "-h" {
			fmt.Fprintln(
				flag.CommandLine.Output(),
				HelpString,
			)
		}
		fmt.Fprintf(
			flag.CommandLine.Output(),
			"Usage:\n  estats [OPTION...]\nOptions:\n",
		)
		flag.PrintDefaults()
	}

	interval := flag.Int(
		"i",
		5,
		"interval in seconds at which to poll",
	)
	path := flag.String(
		"p",
		"",
		"partition to poll",
	)
	net := flag.String(
		"n",
		"",
		"network interface to poll",
	)

	flag.Parse()

	if *interval <= 0 {
		fmt.Fprintf(
			os.Stderr,
			"invalid value \"%d\" for flag -i: out of scope\n",
			*interval,
		)
		flag.PrintDefaults()
		os.Exit(2)
	}

	providers := []stats.Provider{
		stats.LoadAvgProvider(),
		stats.CpuUsageProvider(),
		stats.MemUsageProvider(),
	}

	if *net != "" {
		providers = append(providers, stats.NetUsageProvider(*net))
	}

	if *path != "" {
		providers = append(providers, stats.DiskUsageProvider(*path))
	}

	p := NewPoller(providers...)
	p.Run(*interval)
}
