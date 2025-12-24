package proxy

import (
	"crypto/tls"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ProxyRequest forwards requests to the WHM server
func ProxyRequest(c *gin.Context, whmURL string) {
	// Create new request
	req, err := http.NewRequest(c.Request.Method, whmURL, c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating request")
		return
	}

	// Copy headers
	for key, values := range c.Request.Header {
		for _, v := range values {
			req.Header.Add(key, v)
		}
	}

	// HTTP client (skip SSL verification if needed)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Perform request
	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusBadGateway, "Error connecting to WHM")
		return
	}
	defer resp.Body.Close()

	// Copy response headers & rewrite cookies
	for key, values := range resp.Header {
		if key == "Set-Cookie" {
			for _, v := range values {
				v = strings.ReplaceAll(v, "imatech-taguig.net", c.Request.Host)
				c.Writer.Header().Add(key, v)
			}
		} else {
			for _, v := range values {
				c.Writer.Header().Add(key, v)
			}
		}
	}

	// Copy status code
	c.Status(resp.StatusCode)

	// Rewrite body if HTML
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/html") {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := strings.ReplaceAll(string(bodyBytes), "imatech-taguig.net", c.Request.Host)
		c.Writer.Write([]byte(bodyString))
	} else {
		io.Copy(c.Writer, resp.Body)
	}
}
