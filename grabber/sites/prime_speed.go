package sites

import (
	"proxy_grabber/grabber"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type PrimeSpeed struct {
	url       string
	proxyType grabber.ProxyType
	sync.WaitGroup
}

func NewPrimeSpeed() grabber.Grabber {
	genericRegexp := &GenericRegexp{
		proxyType: grabber.HTTP,
		url:       "http://www.prime-speed.ru/proxy/free-proxy-list/all-working-proxies.php",
		textMatcher: func(doc *goquery.Document) string {
			return doc.Find("pre").First().Text()
		}}

	return genericRegexp
}
