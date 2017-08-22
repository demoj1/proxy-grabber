package proxy_http_client

import (
	"net"
	"net/http"
	"net/url"
	"time"
)

var transport = &http.Transport{Dial: (&net.Dialer{
	Timeout:   30 * time.Second,
	KeepAlive: 30 * time.Second,
}).Dial,
	TLSHandshakeTimeout:   10 * time.Second,
	ResponseHeaderTimeout: 10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

// GetClient - возвращает http клиента с установленым прокси.
// С помощью данного клиента рекомендуется делать только один запрос.
func GetClient(address string) http.Client {
	proxyURL := getProxy(address)
	transport.Proxy = http.ProxyURL(proxyURL)

	return http.Client{
		Timeout:   5 * time.Second,
		Transport: transport}
}

func getProxy(address string) *url.URL {
	URL, err := url.Parse(address)

	if err != nil {
		panic(err)
	}

	return URL
}
