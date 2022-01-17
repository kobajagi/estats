package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/kobajagi/estats/stats"
)

type Result struct {
	data map[string]string
	err  error
}

type Poller struct {
	providers []stats.Provider
}

func NewPoller(p ...stats.Provider) *Poller {
	return &Poller{}
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
		output := map[string]string{}

		for result := range c {
			if result.err != nil {
				fmt.Fprintln(os.Stderr, result.err)
				continue
			}

			for k, v := range result.data {
				output[k] = v
			}
		}

		fmt.Println(output)
	}(outputChan)

	// concurent execution of stats providers
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
