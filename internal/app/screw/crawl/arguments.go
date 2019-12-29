package crawl

import (
	"errors"
	"net/url"
)

func normalize(source []string) ([]*url.URL, error) {
	if len(source) == 0 {
		return nil, errors.New("no arguments")
	}

	var urls []*url.URL

	for _, s := range source {
		u, err := url.Parse(s)
		if err != nil {
			return []*url.URL{}, err
		}
		urls = append(urls, u)
	}

	return urls, nil
}
