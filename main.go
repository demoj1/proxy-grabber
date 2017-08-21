package main

import (
	"log"
	"proxy_grabber/grabber"
	"proxy_grabber/grabber/sites"
	"runtime"
	"time"

	"net/http"

	_ "net/http/pprof"
)

func main() {
	grabber.Registry.Add(
		"fresh", sites.NewFreshProxy(),
	).Add(
		"hidemy", sites.NewHidemy(),
	).Add(
		"multiproxy", sites.NewMultiProxy(),
	).Add(
		"primespeed", sites.NewPrimeSpeed(),
	).Add(
		"therealist", sites.NewThereAList(),
	)

	c, err := grabber.Registry.Grab(grabber.HTTP)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			log.Println("Count active goroutine: ", runtime.NumGoroutine())
			time.Sleep(500 * time.Millisecond)
		}
	}()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	for proxy := range c {
		log.Printf("%v - \x1b[32;1mALIVE\x1b[37;1m", proxy)
	}
}
