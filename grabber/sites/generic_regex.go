package sites

import (
	"fmt"
	"proxy_grabber/grabber"

	"regexp"

	"github.com/PuerkitoBio/goquery"
)

const REGEX = `(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})[ \t:]+(\d{2,5})`

type TextMatcher func(document *goquery.Document) string

type GenericRegexp struct {
	url         string
	proxyType   grabber.ProxyType
	textMatcher TextMatcher
}

func (g *GenericRegexp) Grab(proxyType grabber.ProxyType) error {
	g.proxyType = proxyType

	doc, err := goquery.NewDocument(fmt.Sprintf(g.url))
	if err != nil {
		return err
	}

	text := g.textMatcher(doc)
	re := regexp.MustCompile(REGEX)

	matches := re.FindAllString(text, -1)

	for _, match := range matches {
		grabber.PoolInstance.Tasks <- grabber.Task{
			match,
			g.proxyType}
	}

	return nil
}
