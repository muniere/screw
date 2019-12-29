package spider

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/xmlpath.v2"
)

type IndexOptions struct {
	Focus Focus
	Grep  *regexp.Regexp
}

func Index(uri *url.URL, options IndexOptions) ([]*url.URL, error) {
	switch options.Focus {
	case Href:
		return index(uri, "//a/@href", nil, options)
	case HrefText:
		return index(uri, "//a/@href", regexp.MustCompile(".*\\.txt"), options)
	case HrefImage:
		return index(uri, "//a/@href", regexp.MustCompile(".*\\.(jpg|png|gif)"), options)
	case Image:
		return index(uri, "//img/@src", regexp.MustCompile(".*\\.(jpg|png|gif)"), options)
	case Script:
		return index(uri, "//script/@src", regexp.MustCompile(".*\\.(js)"), options)
	default:
		return []*url.URL{}, nil
	}
}

func index(uri *url.URL, path string, pattern *regexp.Regexp, options IndexOptions) ([]*url.URL, error) {
	log.Debugf("Get contents of URI: %s", uri.String())

	res, err := http.Get(uri.String())
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != 200 {
		return nil, errors.New(res.Status)
	}

	log.Debugf("Search URIs from URI: %s", uri.String())
	doc, err := xmlpath.ParseHTML(res.Body)
	if err != nil {
		return nil, err
	}

	var found []*url.URL

	xpath := xmlpath.MustCompile(path)
	iter := xpath.Iter(doc)

	for iter.Next() {
		val := iter.Node().String()

		if pattern != nil && !pattern.MatchString(val) {
			continue
		}
		if options.Grep != nil && !options.Grep.MatchString(val) {
			continue
		}

		s := strings.Replace(val, " ", "+", -1)
		u, err := url.Parse(s)
		if err != nil {
			continue
		}

		found = append(found, u)
	}

	return found, nil
}
