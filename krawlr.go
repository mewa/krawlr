package krawlr

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

type LinkSet map[string]bool

type Krawlr struct {
	links map[string]*LinkSet

	urlsC chan *url.URL
}

func New() *Krawlr {
	return &Krawlr{
		links: map[string]*LinkSet{},

		urlsC: make(chan *url.URL),
	}
}

func (kr *Krawlr) Crawl(addr string) (map[string]*LinkSet, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return map[string]*LinkSet{}, err
	}

	log.Println("scraping", addr)

	err = kr.crawl(u)
	return kr.links, err
}

func (kr *Krawlr) crawl(root *url.URL) error {
	r, err := http.DefaultClient.Do(&http.Request{
		Method: "GET",
		URL:    root,
	})

	if err != nil {
		return err
	}
	defer r.Body.Close()

	return kr.scrape(r.Body)
}

func (kr *Krawlr) scrape(r io.Reader) error {
	return nil
}

func (ls *LinkSet) String() string {
	s := "{"
	i := 0

	for link := range *ls {
		if i == 0 {
			s += link
		} else {
			s += "," + link
		}
		i++
	}
	s += "}"

	return s
}
