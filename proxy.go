package main

import (
    "io"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
)

// proxyRequest forwards requests from Render domain to WHM subdomain
func proxyRequest(c *gin.Context, whmURL string) {
    // Create a new request to WHM
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

    // Perform the request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        c.String(http.StatusBadGateway, "Error connecting to WHM")
        return
    }
    defer resp.Body.Close()

    // Copy response headers and rewrite cookies
    for key, values := range resp.Header {
        if key == "Set-Cookie" {
            for _, v := range values {
                // Rewrite WHM domain in cookies to Render domain
                v = strings.ReplaceAll(v, "whmdomain.com", c.Request.Host)
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

    // Optional: rewrite absolute URLs in HTML body
    contentType := resp.Header.Get("Content-Type")
    if strings.Contains(contentType, "text/html") {
        bodyBytes, _ := io.ReadAll(resp.Body)
        bodyString := strings.ReplaceAll(string(bodyBytes), "whmdomain.com", c.Request.Host)
        c.Writer.Write([]byte(bodyString))
    } else {
        // For non-HTML responses, just copy raw body
        io.Copy(c.Writer, resp.Body)
    }
}
