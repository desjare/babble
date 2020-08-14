
package main

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/url"
	"time"
)




type WWWDb struct {
	db *sql.DB
}

func (db *WWWDb) Open() (err error) {
	db.db, err = sql.Open("mysql", "desjare:desjare@/lmdb")
	if err != nil {
		panic(err)
	}

	err = db.db.Ping()
	if err != nil {
		panic(err)
	}
	return
}

func (db *WWWDb) InsertRequests(urls []*url.URL) (err error){
	var buffer bytes.Buffer

	if len(urls) == 0 {
		return
	}

	buffer.WriteString("INSERT IGNORE INTO `requests` (`url`) VALUES ")
	for i, url := range(urls) {

		// ignore fragment
		url.Fragment = ""

		buffer.WriteString(fmt.Sprintf("(\"%s\")", url.String()))
		if i != len(urls) - 1 {
			buffer.WriteString(",")
		}
	}

	_, err = db.db.Exec(buffer.String())
	if err == nil {
		return err
	}
	return
}

func (db *WWWDb) FetchRequests(numrequests int) (urls []string, err error) {
	var id int
	var url string

	fetchid := time.Now().UnixNano()

	_, err = db.db.Exec("UPDATE `requests` SET `processed`=1, `fetchid` = ? WHERE `processed` = 0 LIMIT ?", fetchid, numrequests)
	if err != nil {
		return
	}

	rows, err := db.db.Query("SELECT `id`, `url` FROM `requests` WHERE `fetchid` = ?", fetchid)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id, &url)
		if err != nil {
			return
		}
		urls = append(urls, url)
	}
	return
}


