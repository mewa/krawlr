package krawlr

import (
	"errors"
)

var (
	ErrMissingHrefAttr   = errors.New("link is missing href attribute")
	ErrUnsupportedScheme = errors.New("unsupported scheme")
)
