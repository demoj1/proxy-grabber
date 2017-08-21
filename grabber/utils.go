package grabber

import (
	"fmt"
	"proxy_grabber/proxy_http_client"
	"strings"
	"sync"
)

func CheckAddress(address string, proxyType ProxyType) bool {
	switch proxyType {
	case HTTP, HTTPS:
		client := proxy_http_client.GetClient(fmt.Sprintf(
			"%v://%v",
			strings.ToLower(proxyType.String()),
			address))

		resp, err := client.Head("http://blank.org/")
		defer func() {
			if resp != nil {
				resp.Body.Close()
			}
		}()
		if err != nil {
			return false
		}

		if s := resp.StatusCode / 100; s == 2 || s == 3 {
			return true
		}
	}

	return false
}

func Merge(cs ...chan string) chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	output := func(c <-chan string) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
