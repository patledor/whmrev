package proxy

import (
	"crypto/tls"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ProxyRequest(c *gin.Context, whmURL string) {
	req, err := http.NewRequest(c.Request.Method, whmURL, c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Request error")
		return
	}

	for k, v := range c.Request.Header {
		req.Header[k] = v
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusBadGateway, "WHM unreachable")
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		for _, vv := range v {
			if k == "Set-Cookie" {
				vv = strings.ReplaceAll(vv, "imatech-taguig.net", c.Request.Host)
			}
			c.Writer.Header().Add(k, vv)
		}
	}

	c.Status(resp.StatusCode)

	if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		body, _ := io.ReadAll(resp.Body)
		bodyStr := strings.ReplaceAll(string(body), "imatech-taguig.net", c.Request.Host)
		c.Writer.Write([]byte(bodyStr))
	} else {
		io.Copy(c.Writer, resp.Body)
	}
}
