package krawlr

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"net/http"
	"net/url"
)

type LinkSet map[string]bool

type Krawlr struct {
	links map[string]*LinkSet

	client *http.Client

	urlsC chan *url.URL
}

func New() *Krawlr {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return &Krawlr{
		links: map[string]*LinkSet{},

		client: client,

		urlsC: make(chan *url.URL),
	}
}

func (kr *Krawlr) Crawl(addr string) (map[string]*LinkSet, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return map[string]*LinkSet{}, err
	}

	log.WithField("root", u).WithField("url", addr).Println("scraping")

	err = kr.crawl(u, u)
	return kr.links, err
}

func (kr *Krawlr) crawl(root, addr *url.URL) error {
	r, err := kr.client.Do(&http.Request{
		Method: "GET",
		URL:    addr,
	})

	if err != nil {
		return err
	}
	defer r.Body.Close()

	switch r.StatusCode {
	case 301, 302:
		location := r.Header.Get("Location")
		locUrl, _ := url.Parse(location)

		crawled := kr.visitLink(root, locUrl)

		log.
			WithField("url", addr).
			WithField("location", locUrl).
			Println("redirect")

		if !crawled {
			if locUrl.Host == root.Host {
				return kr.crawl(root, locUrl)
			}
		}
		return nil
	case 404:
	default:
		return kr.scrape(root, r.Body)
	}

	return nil
}

func (kr *Krawlr) scrape(root *url.URL, r io.Reader) error {
	tokenizer := html.NewTokenizer(r)

	for {
		tt := tokenizer.Next()
		switch tt {

		case html.ErrorToken:
			err := tokenizer.Err()
			if err != io.EOF {
				log.WithField("url", root).Warn("error token encountered", err)
				return err
			} else {
				return nil
			}

		case html.StartTagToken, html.SelfClosingTagToken:
			tok := tokenizer.Token()
			if tok.DataAtom == atom.A {
				link, err := extractUrl(&tok)
				if err != nil {
					log.
						WithField("root", root).
						WithField("tag", tok).
						Warn("invalid link encountered")
					continue
				}

				err = kr.analyse(root, absUrl(root, link))
				if err != nil {
					log.WithError(err).WithField("link", link).Warn("error analysing link")
				}
			}
		}
	}

	return nil
}

func (kr *Krawlr) analyse(root, link *url.URL) error {
	if len(link.Scheme) != 0 && link.Scheme != "http" && link.Scheme != "https" {
		return ErrUnsupportedScheme
	}

	crawled := kr.visitLink(root, link)

	if !crawled {
		if link.Host == root.Host {
			log.WithField("root", root).WithField("url", link).Println("scraping")

			err := kr.crawl(link, link)
			if err != nil {
				log.WithError(err).Error("error crawling")
			}
		}
	}

	return nil
}

// marks link as visited on page root and returns whether it had been
// visited previously
func (kr *Krawlr) visitLink(root, link *url.URL) bool {
	rootStr := root.String()
	set, ok := kr.links[rootStr]

	if !ok {
		// initialise link set
		set = &LinkSet{}
		kr.links[rootStr] = set
	}

	_, crawled := kr.links[link.String()]
	(*set)[link.String()] = true

	return crawled
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
