package main

import (
	"log"
	"proxy_grabber/grabber"
	"proxy_grabber/grabber/sites"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

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

	grabber.Registry.StartLoop(time.Minute * 5)

	r := gin.Default()
	r.POST("/proxy", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"proxy": grabber.Registry.WorkProxy.Pop(),
		})
	})
	r.Run(":8080")
}
