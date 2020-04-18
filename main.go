package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
)

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
}

type Channel struct {
	XMLName xml.Name `xml:"channel"`
	Items   []Item   `xml:"item"`
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

func main() {
	// read from file
	// bt, err := ioutil.ReadFile("./hackernews.xml")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fetch from API
	res, err := http.Get("https://hnrss.org/newest")
	if err != nil {
		log.Fatal(err)
	}

	bt, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	rss := RSS{}
	err = xml.Unmarshal(bt, &rss)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(rss.Channel.Items)
}
