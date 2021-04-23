package web

import (
	"net"

	"server/web/blocker"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"server/log"
	"server/torr"
	"server/version"
	"server/web/api"
	"server/web/auth"
	"server/web/pages"
)

var (
	BTS      = torr.NewBTS()
	waitChan = make(chan error)
)

func Start(port string) {
	log.TLogln("Start TorrServer", version.Version)
	ips := getLocalIps()
	if len(ips) > 0 {
		log.TLogln("IP:", ips)
	}
	err := BTS.Connect()
	if err != nil {
		waitChan <- err
		return
	}
	gin.SetMode(gin.ReleaseMode)

	route := gin.New()
	route.Use(gin.Recovery(), cors.Default(), blocker.Blocker())

	route.GET("/echo", echo)

	routeAuth := auth.SetupAuth(route)
	if routeAuth != nil {
		api.SetupRoute(routeAuth)
		pages.SetupRoute(routeAuth)
	} else {
		api.SetupRoute(&route.RouterGroup)
		pages.SetupRoute(&route.RouterGroup)
	}
	log.TLogln("Start web", port)
	waitChan <- route.Run(":" + port)
}

func Wait() error {
	return <-waitChan
}

func Stop() {
	BTS.Disconnect()
	waitChan <- nil
}

func echo(c *gin.Context) {
	c.String(200, "%v", version.Version)
}

func getLocalIps() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.TLogln("Error get local IPs")
		return nil
	}
	var list []string
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		if i.Flags&net.FlagUp == net.FlagUp {
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}
				if !ip.IsLoopback() {
					list = append(list, ip.String())
				}
			}
		}
	}
	return list
}
