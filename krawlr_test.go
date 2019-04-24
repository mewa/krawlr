package krawlr

import (
	"reflect"
	"testing"
)

func Test_Krawl_Basic(t *testing.T) {
	k := New(100)

	results, err := k.Crawl("http://localhost:8888/")

	expected := map[string]*LinkSet{
		"http://localhost:8888/": &LinkSet{
			"https://google.com":    true,
			"https://wikipedia.org": true,
		},
	}

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if !reflect.DeepEqual(results, expected) {
		t.Fatalf("expected=%v, actual=%v\n", expected, results)
	}
}

func Test_Krawl_Subpages(t *testing.T) {
	k := New(100)

	results, err := k.Crawl("http://localhost:8888/")

	expected := map[string]*LinkSet{
		"http://localhost:8888/": &LinkSet{
			"https://google.com":             true,
			"https://wikipedia.org":          true,
			"http://localhost:8888/sub.html": true,
		},
		"http://localhost:8888/index.html": &LinkSet{
			"https://google.com":             true,
			"https://wikipedia.org":          true,
			"http://localhost:8888/sub.html": true,
		},
		"http://localhost:8888/sub.html": &LinkSet{
			"https://google.com":               true,
			"https://wikipedia.org":            true,
			"http://localhost:8888/":           true,
			"http://localhost:8888/index.html": true,
			"http://localhost:8888/x":          true,
		},
	}

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if !reflect.DeepEqual(results, expected) {
		t.Fatalf("expected=%v, actual=%v\n", expected, results)
	}
}
