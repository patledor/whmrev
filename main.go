package main

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"go-proxy/proxy"
)

func main() {
	r := gin.Default()

	// Catch-all route for wildcard subdomains
	r.Any("/*path", func(c *gin.Context) {
		host := c.Request.Host // e.g., shop.renderdomain.com
		parts := strings.Split(host, ".")
		if len(parts) < 3 { // expecting subdomain.domain.tld
			c.String(400, "Subdomain required")
			return
		}
		subdomain := parts[0]

		// Build WHM URL dynamically
		whmURL := "https://" + subdomain + ".imatech-taguig.net" + c.Param("path")

		// Forward request via proxy
		proxy.ProxyRequest(c, whmURL)
	})

	port := "8080"
	log.Println("Wildcard reverse proxy running on http://localhost:" + port)
	r.Run(":" + port)
}
