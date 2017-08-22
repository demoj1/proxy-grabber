package sites

import (
	"fmt"
	"proxy_grabber/grabber"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type Hidemy struct {
	url       string
	proxyType grabber.ProxyType
}

func NewHidemy() *Hidemy {
	return &Hidemy{
		url:       "https://hidemy.name/ru/proxy-list/?maxtime=1000&anon=234%v#list",
		proxyType: grabber.HTTP,
	}
}

func (h *Hidemy) Grab(addrType grabber.ProxyType) error {
	h.proxyType = addrType

	doc, err := goquery.NewDocument(fmt.Sprintf(h.url, ""))
	if err != nil {
		return err
	}

	pagesCount, err := strconv.Atoi(doc.Find(".proxy__pagination a").Last().Text())
	if err != nil {
		return err
	}

	for i := 0; i <= pagesCount; i++ {
		doc, err = goquery.NewDocument(fmt.Sprintf(h.url, "&start="+strconv.Itoa(i*64)))
		if err != nil {
			panic(err)
		}

		h.getProxyList(doc)
	}

	return nil
}

func (h *Hidemy) getProxyList(doc *goquery.Document) {
	doc.Find("table.proxy__t tbody tr").Each(func(i int, s *goquery.Selection) {
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

		grabber.PoolInstance.Tasks <- grabber.Task{
			address,
			proxyType}
	})
}
