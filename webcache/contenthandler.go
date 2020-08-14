
package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"
)

func handleContent(data Data, request *Request, channels *Channels) {
	startTime := time.Now()

	os.Mkdir(request.output, 0755)
	os.Mkdir(path.Join(request.output,data.Url.Host), 0755)
	path := path.Join(request.output, data.Url.Host, fmt.Sprintf("%x.html", md5.Sum([]byte(data.Url.String()))))
	err := ioutil.WriteFile(path, data.Content, 0666)
	if err != nil {
		panic(err)
	}

	request.contentDuration = time.Since(startTime)
}

