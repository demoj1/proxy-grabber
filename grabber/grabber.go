package grabber

import (
	"fmt"
	"proxy_grabber/proxy_http_client"
	"strings"
)

type ProxyType int

const (
	HTTP ProxyType = iota
	HTTPS
)

type Grabber interface {
	Grab(ProxyType) (chan string, error)
}

func (p ProxyType) String() string {
	switch p {
	case HTTP:
		return "HTTP"
	case HTTPS:
		return "HTTPS"
	}

	return "NONE"
}

func CheckAddress(address string, proxyType ProxyType) bool {
	switch proxyType {
	case HTTP, HTTPS:
		client := proxy_http_client.GetClient(fmt.Sprintf(
			"%v://%v",
			strings.ToLower(proxyType.String()),
			address))

		resp, err := client.Head("http://www.ya.ru")
		if err != nil {
			return false
		}

		if s := resp.StatusCode / 100; s == 2 || s == 3 {
			return true
		}
	}

	return false
}
