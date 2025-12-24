package main

import (
	"encoding/xml"
	"io"
	"net/http"
)

// Root RSS structure
type RSS struct {
	Channel Channel `xml:"channel"`
}

// Channel data
type Channel struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Items []Item `xml:"item"`
}

// RSS item/article
type Item struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// FetchRSS fetches and parses an RSS feed
func FetchRSS(url string) (*RSS, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rss RSS
	if err := xml.Unmarshal(body, &rss); err != nil {
		return nil, err
	}

	return &rss, nil
}
