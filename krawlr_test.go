package krawlr

import (
	"reflect"
	"testing"
)

func Test_Basic(t *testing.T) {
	k := New()

	results, err := k.Crawl("http://localhost:8888")

	expected := map[string]*LinkSet{
		"http://localhost:8888": &LinkSet{
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
