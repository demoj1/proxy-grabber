package proxy_http_client

import (
	"net/http"
	"net/url"
	"time"
)

// GetClient - возвращает http клиента с установленым прокси.
// С помощью данного клиента рекомендуется делать только один запрос.
func GetClient(address string) *http.Client {
	proxyURL := getProxy(address)

	return &http.Client{
		Timeout: time.Duration(10 * time.Second),
		Transport: &http.Transport{
			Proxy: http.ProxyURL(
				proxyURL)}}
}

func getProxy(address string) *url.URL {
	URL, err := url.Parse(address)

	if err != nil {
		panic(err)
	}

	return URL
}
