package grabber

import (
	"sync"

	"time"

	"github.com/golang-collections/collections/stack"
)

type registry struct {
	grabbers map[string]Grabber

	WorkListMutex *sync.Mutex
	WorkProxy     *stack.Stack

	ListMutex *sync.Mutex
	ProxyList []Proxy
}

const (
	STACK_MAX_SIZE = 500
)

var Registry *registry = &registry{
	grabbers: make(map[string]Grabber),

	WorkListMutex: new(sync.Mutex),
	WorkProxy:     new(stack.Stack),

	ListMutex: new(sync.Mutex),
	ProxyList: make([]Proxy, 0)}

func (r *registry) Add(name string, grabber Grabber) *registry {
	r.grabbers[name] = grabber
	return r
}

func (r *registry) Delete(name string) *registry {
	delete(r.grabbers, name)
	return r
}

func (r *registry) StartLoop(updateInterval time.Duration) {
	r.updateProxyList(updateInterval)

	go func(registry *registry) {
		page := 1
		wg := sync.WaitGroup{}
		for {
			proxy := paginateList(registry, page)
			r.filtering(&wg, proxy)
			r.clearStack()
		}
	}(r)
}
func (r *registry) clearStack() {
	r.WorkListMutex.Lock()
	if r.WorkProxy.Len() > STACK_MAX_SIZE {
		var tmpProxy []Proxy
		for i := 0; i < 10; i++ {
			tmpProxy = append(tmpProxy, r.WorkProxy.Pop().(Proxy))
		}
		r.WorkProxy = new(stack.Stack)
		for _, proxy := range tmpProxy {
			r.WorkProxy.Push(proxy)
		}
	}
	r.WorkListMutex.Unlock()
}
func (r *registry) filtering(wg *sync.WaitGroup, proxy []Proxy) {
	wg.Add(len(proxy))
	for _, p := range proxy {
		go r.task(p, wg)
	}
	wg.Wait()
}
func paginateList(registry *registry, page int) []Proxy {
	var proxy []Proxy
	if len(registry.ProxyList) < page*STACK_MAX_SIZE/2 {
		proxy = registry.ProxyList
		page = 1
	} else {
		proxy = registry.ProxyList[page : page+STACK_MAX_SIZE/2]
	}
	return proxy
}

func (r *registry) task(proxy Proxy, wg *sync.WaitGroup) {
	if CheckAddress(proxy) {
		r.WorkListMutex.Lock()
		proxy.Alive = true
		r.WorkProxy.Push(proxy)
		r.WorkListMutex.Unlock()
	}

	wg.Done()
}

func (r *registry) updateProxyList(updateInterval time.Duration) {
	updateProxyListTicker := time.NewTicker(updateInterval)

	updateFun := func() {
		var resList []Proxy
		for _, grabber := range r.grabbers {
			_, list := grabber.Grab()
			resList = append(resList, list...)
		}

		r.ListMutex.Lock()
		resList = removeDuplicates(resList)[:500]
		r.ProxyList = resList
		r.ListMutex.Unlock()
	}

	go updateFun()

	go func() {
		for range updateProxyListTicker.C {
			updateFun()
		}
	}()
}

func removeDuplicates(elements []Proxy) []Proxy {
	encountered := map[Proxy]bool{}
	result := []Proxy{}

	for v := range elements {
		if encountered[elements[v]] == true {
			continue
		} else {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}

	return result
}
