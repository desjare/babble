
package main

import (
	"mucus/babble/words"
	"flag"
	"net/url"
	"fmt"
	"sync"
	"time"
)

const (
	MaxRequest = 32
)

type Data struct {
	Url *url.URL
	Content string
	Links map[string][]*url.URL
}

type Request struct {
	url url.URL
	context *words.TokenizeContext

	// profiling
	browseDuration time.Duration
	urlDuration time.Duration
	tokenizeDuration time.Duration
	contentDuration time.Duration

	// stats
	numWords int
	numURLs int
}

type Channels struct {
	quitch chan int

	db LMDb
	dblock sync.Mutex

	bench bool
}

func browseHandler(data Data, request *Request, channels *Channels) {
	handleURL(data, request, channels)
	handleContent(data, request, channels)
}


func browseURL(rawurl string, bench bool) {
	var channels Channels

	u := geturl(rawurl)

	channels.bench = bench
	err := channels.db.Open()
	if err != nil {
		panic(err)
	}

	context, err := words.TokenizeNewContext()
	if err != nil {
		panic("cannot initialize tokenizer")
	}

	// browse initial url
	request := Request{*u, context, 0, 0, 0, 0, 0, 0}
	go browse(request, &channels, browseHandler)

	for {
		// dispatch requests
		var urls []string

		urls, err = channels.db.FetchRequests(MaxRequest)
		if err != nil {
			panic(err)
		}

		if len(urls) == 0 {
			fmt.Printf("No request. Waiting.\n")
			time.Sleep(1000 * time.Millisecond)
			continue
		}

		var wg sync.WaitGroup
		startTime := time.Now()
		for _, s := range(urls) {

			u, err := url.Parse(s)
			if err != nil {
				continue
			}

			request := Request{*u, context, 0, 0, 0, 0, 0, 0}

			wg.Add(1)
			go func(request Request, channels *Channels, browseHandler func(data Data, request *Request, channels *Channels) ) {
				browse(request, channels, browseHandler)
				wg.Done()
			} (request, &channels, browseHandler)
		}
		wg.Wait()
		fmt.Printf("Processing %d urls done in %s\n", len(urls), time.Since(startTime).String() )
	}
}


func buildlm(path string) {
	dict := words.Dictionary{}
	err := dict.ReadXML()
	if err != nil {
		panic(err)
	}
	err = dict.WriteBinary(path)
	if err != nil {
		panic(err)
	}
}


func main() {

	var rawurl string
	var lmpath string
	var bench bool

	flag.StringVar(&rawurl, "url", "", "url to search")
	flag.StringVar(&lmpath, "buildlm", "", "build lm")
	flag.BoolVar(&bench, "bench", false, "output benchmarks")
	flag.Parse()

	if len(rawurl) > 0 {
		browseURL(rawurl, bench)
	}

	if len(lmpath) > 0 {
		buildlm(lmpath)
	}

}
