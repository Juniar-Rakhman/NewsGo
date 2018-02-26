package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type SiteMapIndex struct {
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

type NewsAggPage struct {
	Title string
	News  map[string]NewsMap
}

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Go index")
}

func newsAggHandler(w http.ResponseWriter, _ *http.Request) {
	var s SiteMapIndex
	var n News

	resp, _ := http.Get("https://www.washingtonpost.com/news-sitemap-index.xml")
	bytes, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &s)
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

	p := NewsAggPage{Title: "Amazing News Aggregator", News: newsMap}
	t, _ := template.ParseFiles("web/news.html")
	t.Execute(w, p)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/agg/", newsAggHandler)
	http.ListenAndServe(":8000", nil)
}
