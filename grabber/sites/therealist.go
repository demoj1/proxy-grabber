package sites

import (
	"proxy_grabber/grabber"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type ThereAList struct {
	url       string
	proxyType grabber.ProxyType
	sync.WaitGroup
}

func NewThereAList() grabber.Grabber {
	genericRegexp := &GenericRegexp{
		proxyType: grabber.HTTP,
		url:       "http://www.therealist.ru/proksi/spisok-vsex-rabochix-proksi",
		textMatcher: func(doc *goquery.Document) string {
			return doc.Find("pre").First().Text()
		}}

	return genericRegexp
}
