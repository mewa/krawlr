package krawlr

import (
	"golang.org/x/net/html"
	"net/url"
	"path"
)

func href(attrs []html.Attribute) *html.Attribute {
	for _, attr := range attrs {
		if attr.Key == "href" {
			return &attr
		}
	}
	return nil
}

// extracts url from <a> tag
func extractUrl(tok *html.Token) (*url.URL, error) {
	h := href(tok.Attr)

	if h != nil {
		link, err := url.Parse(h.Val)

		if err != nil {
			return nil, err
		}

		return link, nil
	}
	return nil, ErrMissingHrefAttr
}

func absUrl(root, link *url.URL) *url.URL {
	siteRelative := len(link.Host) == 0

	if siteRelative {
		if path.IsAbs(link.Path) {
			return root.ResolveReference(link)
		} else {
			ret := *root

			ret.Path = path.Join(ret.Path, link.Path)
			return &ret
		}
	}

	return link
}
