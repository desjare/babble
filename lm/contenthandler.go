
package main

import (
	"mucus/babble/words"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func handleContent(data Data, request *Request, channels *Channels) {
	startTime := time.Now()
	tokens := words.Tokenize(data.Content, request.context, true)
	wordmap := make(map[string] int)
	request.tokenizeDuration = time.Since(startTime)

	for _, tok := range tokens {
		if tok.Word != nil && !tok.Word.IsPunct() {
			wordmap[tok.Word.String()]++
		}
	}

	bytes, err := json.Marshal(wordmap)
	if err != nil {
		panic(err)
	}

	os.Mkdir("lmoutput", 0755)
	path := "lmoutput/" + fmt.Sprintf("%x.json", md5.Sum([]byte(data.Url.String())))
	err = ioutil.WriteFile(path, bytes, 0666)
	if err != nil {
		panic(err)
	}

	request.contentDuration = time.Since(startTime)
	request.numWords = len(wordmap)
}

