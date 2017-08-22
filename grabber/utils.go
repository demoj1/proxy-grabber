package grabber

import (
	"fmt"
	"io"
	"io/ioutil"
	"proxy_grabber/proxy_http_client"
	"strings"
)

var (
	PoolInstance *Pool
)

type Task struct {
	Proxy     string
	ProxyType ProxyType
}

type Pool struct {
	Workers int

	Results chan string
	Tasks   chan Task

	start bool
}

func InitPool(workers int) *Pool {
	if PoolInstance == nil {
		PoolInstance = &Pool{
			Workers: workers,
			Results: make(chan string),
			Tasks:   make(chan Task)}
	}

	return PoolInstance
}

func (p *Pool) Start() {
	p.start = true

	for i := 0; i < p.Workers; i++ {
		go p.Do()
	}
}

func (p *Pool) Stop() {
	p.start = false
	close(p.Tasks)
	close(p.Results)
}

func (p *Pool) Do() {
	for {
		if !p.start {
			return
		}

		task := <-p.Tasks
		if checkAddress(task.Proxy, task.ProxyType) {
			p.Results <- task.Proxy
		}
	}
}

func checkAddress(address string, proxyType ProxyType) bool {
	switch proxyType {
	case HTTP, HTTPS:
		client := proxy_http_client.GetClient(fmt.Sprintf(
			"%v://%v",
			strings.ToLower(proxyType.String()),
			address))

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
