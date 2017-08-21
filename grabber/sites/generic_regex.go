package sites

import (
	"fmt"
	"proxy_grabber/grabber"
	"sync"

	"regexp"

	"github.com/PuerkitoBio/goquery"
)

const REGEX = `(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})[ \t:]+(\d{2,5})`

type TextMatcher func(document *goquery.Document) string

type GenericRegexp struct {
	url         string
	proxyType   grabber.ProxyType
	textMatcher TextMatcher

	sync.WaitGroup
}

func (g *GenericRegexp) Grab(proxyType grabber.ProxyType) (chan string, error) {
	outchan := make(chan string, 5)
	g.proxyType = proxyType

	doc, err := goquery.NewDocument(fmt.Sprintf(g.url))
	if err != nil {
		return nil, err
	}

	text := g.textMatcher(doc)
	re := regexp.MustCompile(REGEX)

	matches := re.FindAllString(text, -1)

	g.Add(len(matches))
	for _, match := range matches {
		go func(proxy string) {
			if grabber.CheckAddress(proxy, g.proxyType) {
				outchan <- proxy
			}
			g.Done()
		}(match)
	}

	go func() {
		g.Wait()
		close(outchan)
	}()

	return outchan, nil
}
