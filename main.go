package main

import (
	"log"
	"proxy_grabber/grabber"
	"proxy_grabber/grabber/sites"
)

func main() {
	a := sites.NewPrimeSpeed()
	c, _ := a.Grab(grabber.HTTP)

	for proxy := range c {
		log.Printf("%v - \x1b[32;1mALIVE\x1b[37;1m", proxy)
	}
}

//"http://www.prime-speed.ru/proxy/free-proxy-list/all-working-proxies.php", [+]
//"http://fineproxy.org/freshproxy", [+]

//"http://www.therealist.ru/proksi/spisok-vsex-rabochix-proksi",
//"http://multiproxy.org/txt_anon/proxy.txt",
//"http://spys.ru/en/http-proxy-list/" + POST DATA {
//	xpp:5
//	xf5:1
//}
