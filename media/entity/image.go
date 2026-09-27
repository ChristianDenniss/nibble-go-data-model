package entity

import (
	"errors"
	"net/url"
	"strings"
)

var ErrImageURLInvalid = errors.New("image url must be an absolute http(s) url")

// NormalizeImageURL trims raw and checks it is an absolute http(s) URL. Empty means "no image".
// Images are external links for now; a Nibble-hosted bucket later is still just a URL here.
func NormalizeImageURL(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", nil
	}
	u, err := url.Parse(trimmed)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return "", ErrImageURLInvalid
	}
	return trimmed, nil
}
