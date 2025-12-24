package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// proxyRequest forwards requests from Render domain to WHM subdomain
func proxyRequest(c *gin.Context, whmURL string) {
	// Create new request to WHM
	req, err := http.NewRequest(c.Request.Method, whmURL, c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating request")
		return
	}

	// Copy headers from original request
	for key, values := range c.Request.Header {
		for _, v := range values {
			req.Header.Add(key, v)
		}
	}

	// HTTP client (skip TLS verify if needed)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // remove in production if using valid certs
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusBadGateway, "Error connecting to WHM")
		return
	}
	defer resp.Body.Close()

	// Copy response headers and rewrite cookies
	for key, values := range resp.Header {
		for _, v := range values {
			if key == "Set-Cookie" {
				v = strings.ReplaceAll(v, "imatech-taguig.net", c.Request.Host)
			}
			c.Writer.Header().Add(key, v)
		}
	}

	// Copy status code
	c.Status(resp.StatusCode)

	// Copy body
	if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := strings.ReplaceAll(string(bodyBytes), "imatech-taguig.net", c.Request.Host)
		c.Writer.Write([]byte(bodyString))
	} else {
		io.Copy(c.Writer, resp.Body)
	}
}

func main() {
	r := gin.Default()

	// Catch-all route for wildcard subdomains
	r.Any("/*path", func(c *gin.Context) {
		host := c.Request.Host
		parts := strings.Split(host, ".")
		if len(parts) < 3 {
			c.String(http.StatusBadRequest, "Subdomain required")
			return
		}
		subdomain := parts[0] // e.g., 'shop' from 'shop.renderdomain.com'

		whmURL := "https://" + subdomain + ".imatech-taguig.net" + c.Param("path")

		proxyRequest(c, whmURL)
	})

	port := "8080"
	log.Println("Wildcard reverse proxy running on http://localhost:" + port)
	r.Run(":" + port)
}
