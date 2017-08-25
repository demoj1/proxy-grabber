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
	out       chan string
}

func NewHidemy() *Hidemy {
	return &Hidemy{
		url: "https://hidemy.name/ru/proxy-list/?maxtime=1000&anon=234%v#list",
		out: make(chan string),
	}
}

func (h *Hidemy) Grab() (error, []grabber.Proxy) {
	doc, err := goquery.NewDocument(fmt.Sprintf(h.url, ""))
	if err != nil {
		return err, nil
	}

	pagesCount, err := strconv.Atoi(doc.Find(".proxy__pagination a").Last().Text())
	if err != nil {
		return err, nil
	}

	var proxyList []grabber.Proxy
	for i := 0; i <= pagesCount; i++ {
		doc, err = goquery.NewDocument(fmt.Sprintf(h.url, "&start="+strconv.Itoa(i*64)))
		if err != nil {
			panic(err)
		}

		list := h.getProxyList(doc)
		proxyList = append(proxyList, list...)
	}

	return nil, proxyList
}

func (h *Hidemy) getProxyList(doc *goquery.Document) []grabber.Proxy {
	var proxyList []grabber.Proxy

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

		proxyList = append(proxyList, grabber.Proxy{
			Address: address,
			Type:    proxyType})
	})

	return proxyList
}

func (h *Hidemy) ProxyChan() chan string {
	return h.out
}
