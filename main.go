package main

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/patledor/whmrev/proxy"
)

func main() {
	r := gin.Default()

	r.Any("/*path", func(c *gin.Context) {
		host := c.Request.Host
		parts := strings.Split(host, ".")
		if len(parts) < 3 {
			c.String(400, "Subdomain required")
			return
		}

		subdomain := parts[0]
		whmURL := "https://" + subdomain + ".imatech-taguig.net" + c.Param("path")

		proxy.ProxyRequest(c, whmURL)
	})

	log.Println("Proxy running on :8080")
	r.Run(":8080")
}
