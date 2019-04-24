package krawlr

import (
	"net/url"
	"testing"
)

func Test_Utils_Urls(t *testing.T) {
	root, _ := url.Parse("https://google.com/first")

	var u *url.URL

	cases := [][]string{
		[]string{"/ymca", "https://google.com/ymca"},
		[]string{"./ymca", "https://google.com/first/ymca"},
		[]string{"abba", "https://google.com/first/abba"},
		[]string{"../second", "https://google.com/second"},
		[]string{"https://google.com/third", "https://google.com/third"},
		[]string{"https://wikipedia.org/page", "https://wikipedia.org/page"},
		[]string{"./", "https://google.com/first"},
	}

	for no, kase := range cases {
		u2, _ := url.Parse(kase[0])
		u = absUrl(root, u2)

		var expected string = kase[1]
		if u.String() != expected {
			t.Errorf("case %d (%s): expected=%s, actual=%s", no, kase[0], kase[1], u.String())
		}
	}
}
