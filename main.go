package main

import (
    "log"
    "strings"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // Catch-all route
    r.Any("/*path", func(c *gin.Context) {
        host := c.Request.Host // e.g., shop.renderdomain.com
        var whmURL string

        // Map subdomains to WHM
        switch {
        case strings.HasPrefix(host, "shop."):
            whmURL = "https://shop.whmdomain.com" + c.Param("path")
        case strings.HasPrefix(host, "blog."):
            whmURL = "https://blog.whmdomain.com" + c.Param("path")
        default:
            whmURL = "https://www.whmdomain.com" + c.Param("path")
        }

        proxyRequest(c, whmURL)
    })

    port := "8080"
    log.Println("Reverse proxy running on http://localhost:" + port)
    r.Run(":" + port)
}
