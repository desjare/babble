package main

import (
	"fmt"
	"time"
)

func handleURL(data Data, request *Request, channels *Channels) {

	startTime := time.Now()
	channels.dblock.Lock()
	defer channels.dblock.Unlock()
	err := channels.db.InsertRequests(data.Links[request.url.Host])
	if err != nil {
		fmt.Printf("Error inserting urls %s", err)
	}
	request.urlDuration = time.Since(startTime)
	request.numURLs = len(data.Links[request.url.Host])
}
