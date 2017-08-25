package sites

import (
	"proxy_grabber/grabber"

	"github.com/PuerkitoBio/goquery"
)

func NewFreshProxy() grabber.Grabber {
	genericRegexp := &GenericRegexp{
		url:       "http://fineproxy.org/freshproxy",
		textMatcher: func(doc *goquery.Document) string {
			return doc.Find(".entry-content p").First().Text()
		},
		out: make(chan string)}

	return genericRegexp
}
