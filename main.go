package main

import (
	"log"
	"proxy_grabber/grabber/hidemy"
)

func main() {
	a := hidemy.New()
	c, _ := a.Grab(1)

	for proxy := range c {
		log.Printf("%v - \x1b[32;1mALIVE\x1b[37;1m", proxy)
	}
}
