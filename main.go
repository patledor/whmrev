package main

import (
    "log"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // Catch-all route
    r.Any("/*path", func(c *gin.Context) {
        host := c.Request.Host // e.g., shop.yourdomain.com
        subdomain := strings.Split(host, ".")[0] // get 'shop' from 'shop.yourdomain.com'

        whmURL := "https://" + subdomain + ".imatech-taguig.net" + c.Param("path")

        proxyRequest(c, whmURL)
    })

    port := "8080"
    log.Println("Wildcard reverse proxy running on http://localhost:" + port)
    r.Run(":" + port)
}
