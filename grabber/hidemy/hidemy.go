package hidemy

import (
	"fmt"
	"proxy_grabber/grabber"
	"proxy_grabber/proxy_http_client"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type Hidemy struct {
	url      string
	addrType grabber.ADDR_TYPE
	wg       sync.WaitGroup
}

func New() *Hidemy {
	return &Hidemy{
		url:      "https://hidemy.name/ru/proxy-list/?maxtime=1000&anon=234%v#list",
		addrType: 1,
	}
}

func (h *Hidemy) checkAddress(address, typeNode string) bool {
	types := strings.Split(typeNode, ", ")

	proxy_status := false

	for _, element := range types {
		switch element {
		case "HTTP", "HTTPS":
			client := proxy_http_client.GetClient(strings.ToLower(element) + "://" + address)
			resp, err := client.Head("http://www.ya.ru")
			if err != nil {
				continue
			}

			if s := resp.StatusCode / 100; s == 2 || s == 3 {
				proxy_status = true
			}

		default:
			continue
		}
	}

	return proxy_status
}

func (h *Hidemy) getProxyList(doc *goquery.Document, in chan string) {
	doc.Find("table.proxy__t tbody tr").Each(func(i int, s *goquery.Selection) {
		go func(s *goquery.Selection) {
			h.wg.Add(1)
			tds := s.Find("td").Nodes

			address := tds[0].LastChild.Data + ":" + tds[1].LastChild.Data
			typeNode := tds[4].LastChild.Data

			if h.checkAddress(address, typeNode) {
				in <- address
			}
			h.wg.Done()
		}(s)
	})
}

func (h *Hidemy) Grab(addrType grabber.ADDR_TYPE) (chan string, error) {
	outchan := make(chan string, 50)

	h.addrType = addrType

	doc, err := goquery.NewDocument(fmt.Sprintf(h.url, ""))
	if err != nil {
		return nil, err
	}

	pagesCount, err := strconv.Atoi(doc.Find(".proxy__pagination a").Last().Text())
	if err != nil {
		return nil, err
	}

	h.wg.Add(pagesCount)

	for i := 0; i <= pagesCount; i++ {
		go func(number int) {
			doc, err = goquery.NewDocument(fmt.Sprintf(h.url, "&start="+strconv.Itoa(number*64)))
			if err != nil {
				panic(err)
			}

			h.getProxyList(doc, outchan)
			h.wg.Done()
		}(i)
	}

	go func() {
		h.wg.Wait()
		close(outchan)
	}()

	return outchan, nil
}
