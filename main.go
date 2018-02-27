package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"sync"
)

type siteMap struct {
	Locations []string `xml:"sitemap>loc"`
}

type News struct {
	Titles    []string `xml:"url>news>title"`
	Keywords  []string `xml:"url>news>keywords"`
	Locations []string `xml:"url>loc"`
}

type NewsMap struct {
	Keyword  string
	Location string
}

type NewsPage struct {
	Title string
	News  map[string]NewsMap
}

var wg sync.WaitGroup

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Go index")
}

func newsHandler(w http.ResponseWriter, _ *http.Request) {
	var s siteMap
	var n News

	if resp, err := http.Get("https://www.washingtonpost.com/news-sitemap-index.xml"); err == nil {
		if bytes, err := ioutil.ReadAll(resp.Body); err == nil {
			xml.Unmarshal(bytes, &s)
		}
	} else {
		fmt.Println("shit happens : %s", err)
	}

	newsMap := make(map[string]NewsMap)

	for _, Location := range s.Locations {

		if resp, err := http.Get(Location); err == nil {
			if bytes, err := ioutil.ReadAll(resp.Body); err == nil {
				xml.Unmarshal(bytes, &n)
			}
		} else {
			fmt.Println("shit happens : %s", err)
		}

		for idx := range n.Keywords {
			newsMap[n.Titles[idx]] = NewsMap{n.Keywords[idx], n.Locations[idx]}
		}
	}

	p := NewsPage{Title: "Go News", News: newsMap}
	t, _ := template.ParseFiles("web/news.html")
	t.Execute(w, p)
}

func main() {
	fmt.Println("Serving http://localhost:8000/news/")
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/news/", newsHandler)
	http.ListenAndServe(":8000", nil)
}
