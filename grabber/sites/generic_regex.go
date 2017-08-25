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
	textMatcher TextMatcher
	out         chan string
}

func (g *GenericRegexp) Grab() (error, []grabber.Proxy) {
	doc, err := goquery.NewDocument(fmt.Sprintf(g.url))
	if err != nil {
		return err, nil
	}

	text := g.textMatcher(doc)
	re := regexp.MustCompile(REGEX)

	var res []grabber.Proxy

	matches := re.FindAllString(text, -1)
	for _, match := range matches {
		res = append(res, grabber.Proxy{
			Type:    grabber.HTTP,
			Address: match})
	}
	return nil, res
}
