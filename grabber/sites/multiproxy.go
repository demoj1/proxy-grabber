package sites

import (
	"proxy_grabber/grabber"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type MultiProxy struct {
	url       string
	proxyType grabber.ProxyType
	sync.WaitGroup
}

func NewMultiProxy() grabber.Grabber {
	genericRegexp := &GenericRegexp{
		url: "http://multiproxy.org/txt_anon/proxy.txt",
		textMatcher: func(doc *goquery.Document) string {
			return doc.Text()
		},
		out: make(chan string)}

	return genericRegexp
}
