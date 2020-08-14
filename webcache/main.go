
package main

import (
	"flag"
	"net/url"
	"fmt"
	"sync"
	"time"
	"runtime"
)

const (
	MaxRequest = 64
)

type Data struct {
	Url *url.URL
	Content []byte
	Links map[string][]*url.URL
}

type Request struct {
	url url.URL
	output string

	// profiling
	browseDuration time.Duration
	urlDuration time.Duration
	contentDuration time.Duration

	// stats
	numURLs int
}

type Channels struct {
	db WWWDb
	dblock sync.Mutex

	bench bool
}

func browseHandler(data Data, request *Request, channels *Channels) {
	handleURL(data, request, channels)
	handleContent(data, request, channels)
}


func browseURL(rawurl string, output string, bench bool) {
	var channels Channels

	u := geturl(rawurl)

	channels.bench = bench
	err := channels.db.Open()
	if err != nil {
		panic(err)
	}

	// browse initial url
	request := Request{*u, output, 0, 0, 0, 0}
	go browse(request, &channels, browseHandler)

	for {

		if runtime.NumGoroutine() >= MaxRequest {
			runtime.Gosched()
			continue
		}

		// dispatch requests
		var urls []string

		urls, err = channels.db.FetchRequests(MaxRequest - runtime.NumGoroutine())
		if err != nil {
			panic(err)
		}

		if len(urls) == 0 {
			fmt.Printf("No request. Waiting.\n")
			time.Sleep(1000 * time.Millisecond)
			continue
		} else {
			fmt.Printf("Fetching %d...", len(urls))
		}

		for _, s := range(urls) {

			u, err := url.Parse(s)
			if err != nil {
				continue
			}

			request := Request{*u, output, 0, 0, 0, 0}

			go browse(request, &channels, browseHandler)
		}
	}
}



func main() {

	var rawurl string
	var output string
	var bench bool

	flag.StringVar(&rawurl, "url", "", "url to search")
	flag.StringVar(&output, "output", "www", "output")
	flag.BoolVar(&bench, "bench", false, "output benchmarks")
	flag.Parse()

	if len(rawurl) > 0 {
		browseURL(rawurl, output, bench)
	}
}
