package grabber

import (
	"fmt"
	"io"
	"io/ioutil"
	"proxy_grabber/proxy_http_client"
	"strings"
)

func CheckAddress(proxy Proxy) bool {
	switch proxy.Type {
	case HTTP, HTTPS:
		client := proxy_http_client.GetClient(fmt.Sprintf(
			"%v://%v",
			strings.ToLower(proxy.Type.String()),
			proxy.Address))

		resp, err := client.Get("http://blank.org/")
		if err != nil {
			return false
		}
		defer func() {
			if resp == nil {
				return
			}

			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		}()

		if s := resp.StatusCode / 100; s == 2 || s == 3 {
			return true
		}
	}

	return false
}
