package sites

import (
	"fmt"
	"proxy_grabber/grabber"
	"strconv"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type Hidemy struct {
	url       string
	proxyType grabber.ProxyType
	sync.WaitGroup
}

func NewHidemy() *Hidemy {
	return &Hidemy{
		url:       "https://hidemy.name/ru/proxy-list/?maxtime=1000&anon=234%v#list",
		proxyType: grabber.HTTP,
	}
}

func (h *Hidemy) Grab(addrType grabber.ProxyType) (chan string, error) {
	outchan := make(chan string, 5)

	h.proxyType = addrType

	doc, err := goquery.NewDocument(fmt.Sprintf(h.url, ""))
	if err != nil {
		return nil, err
	}

	pagesCount, err := strconv.Atoi(doc.Find(".proxy__pagination a").Last().Text())
	if err != nil {
		return nil, err
	}
	h.Add(pagesCount)

	for i := 0; i <= pagesCount; i++ {
		go func(number int) {
			doc, err = goquery.NewDocument(fmt.Sprintf(h.url, "&start="+strconv.Itoa(number*64)))
			if err != nil {
				panic(err)
			}

			h.getProxyList(doc, outchan)
			h.Done()
		}(i)
	}

	go func() {
		h.Wait()
		close(outchan)
	}()

	return outchan, nil
}

func (h *Hidemy) getProxyList(doc *goquery.Document, in chan string) {
	doc.Find("table.proxy__t tbody tr").Each(func(i int, s *goquery.Selection) {
		go func(s *goquery.Selection) {
			h.Add(1)
			tds := s.Find("td").Nodes

			address := tds[0].LastChild.Data + ":" + tds[1].LastChild.Data
			typeNodeData := tds[4].LastChild.Data

			var proxyType grabber.ProxyType
			switch typeNodeData {
			case "HTTP", "HTTP, HTTPS":
				proxyType = grabber.HTTP
			case "HTTPS":
				proxyType = grabber.HTTPS
			}

			if proxyType != h.proxyType {
				return
			}

			if grabber.CheckAddress(address, proxyType) {
				in <- address
			}
			h.Done()
		}(s)
	})
}
