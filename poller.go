package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/kobajagi/estats/stats"
)

type Result struct {
	data []stats.Stat
	err  error
}

type Report struct {
	Timestamp int64        `json:"timestamp"`
	Metrics   []stats.Stat `json:"metrics"`
}

type Poller struct {
	providers []stats.Provider
}

func NewPoller(p ...stats.Provider) *Poller {
	return &Poller{providers: p}
}

func (p *Poller) Run(interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)

	for {
		<-ticker.C
		p.Poll()
	}
}

func (p *Poller) Poll() {
	outputChan := make(chan *Result, len(p.providers))
	wg := sync.WaitGroup{}
	wg.Add(len(p.providers))

	// output writer
	go func(c <-chan *Result) {
		report := Report{
			Timestamp: time.Now().UTC().Unix(),
			Metrics:   []stats.Stat{},
		}
		for result := range c {
			if result.err != nil {
				fmt.Fprintln(os.Stderr, result.err)
				continue
			}

			report.Metrics = append(report.Metrics, result.data...)
		}

		fmt.Println(report)
	}(outputChan)

	// concurent execution of metric providers
	for _, provider := range p.providers {
		go func(c chan<- *Result, provider stats.Provider) {
			data, err := provider.Read()
			if err != nil {
				c <- &Result{data: nil, err: err}
			}

			c <- &Result{data: data, err: nil}
			wg.Done()
		}(outputChan, provider)
	}

	wg.Wait()
	close(outputChan)
}
