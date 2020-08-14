package main

import (
	"mucus/babble/words"
        "fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func httpget(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("httpget error get %s\n",err)
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Printf("httpget error read %s\n",err)
	}
	return body
}

func geturl(rawurl string) *url.URL {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil
	}

	// scheme
	if len(u.Scheme) == 0 {
		u.Scheme = "http"
	} else if strings.Compare(u.Scheme, "javascript") == 0 {
		return nil
	} else if strings.Compare(u.Scheme, "mailto") == 0 {
		return nil
	}
	return u
}

func geturls(u *url.URL, links []string) *map[string][]*url.URL {
	urls := make(map[string][]*url.URL)
	for _, link := range links {
		linku := geturl(link)
		if linku != nil {
			if strings.Compare(linku.Host,"") == 0 {
				linku.Host = u.Host
			}
			hosturls := urls[linku.Host]
			urlfound := false
			for _, hosturl := range hosturls {
				if linku.String() == hosturl.String() {
					urlfound = true
					break
				}
			}
			if !urlfound {
				urls[linku.Host] = append(urls[linku.Host], linku)
			}

		}
	}
	return &urls
}

func browse(request Request, channels *Channels, browseHandler func(Data, *Request, *Channels) ) {
	startTime := time.Now()
	body := httpget(request.url.String())
	if body == nil {
		return
	}
	_, links := words.HTMLParse(body)
	urls := geturls(&request.url, links)

	data := Data { &request.url, body, *urls }
	request.browseDuration = time.Since(startTime)
	browseHandler(data, &request, channels)

	if channels.bench {
		fmt.Printf("browse urls %d time: http %s url %s content %s\n",
			request.numURLs, request.browseDuration, request.urlDuration,
			request.contentDuration)
	}
}

