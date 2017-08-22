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
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	pool := grabber.InitPool(10000)
	pool.Start()

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

	go func() {
		for {
			time.Sleep(time.Second * 1)
			log.Println(runtime.NumGoroutine())
		}
	}()

	go grabber.Registry.Grab(grabber.HTTP)

	for proxy := range pool.Results {
		log.Printf("%v - \x1b[32;1mALIVE\x1b[37;1m", proxy)
	}
}
