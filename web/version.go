package web

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lucperkins/rek"

	"github.com/Cloudbox/cvm/build"
)

type versionRequest struct {
	URL string `form:"url" binding:"required"`
}

func (c *Client) Version(g *gin.Context) {
	// parse query
	b := new(versionRequest)
	if err := g.ShouldBindQuery(b); err != nil {
		_ = g.AbortWithError(http.StatusBadRequest, fmt.Errorf("bind query: %w", err))
		return
	}

	// create request
	res, err := rek.Get(b.URL, rek.Timeout(30*time.Second), rek.UserAgent(build.UserAgent))
	if err != nil {
		_ = g.AbortWithError(http.StatusInternalServerError, fmt.Errorf("request url: %w", err))
		return
	}
	defer res.Body().Close()

	// validate response
	if res.StatusCode() != http.StatusOK {
		_ = g.AbortWithError(http.StatusInternalServerError, fmt.Errorf("validate url response: %s", res.Status()))
		return
	}

	contentType := strings.ToLower(res.Headers()["Content-Type"])
	if !strings.Contains(contentType, "xml") && !strings.Contains(contentType, "json") {
		_ = g.AbortWithError(http.StatusInternalServerError, fmt.Errorf("validate url response content-type: %v", contentType))
		return
	}

	// read response
	rb, err := ioutil.ReadAll(http.MaxBytesReader(nil, res.Body(), c.maxResponseSize<<20+1))
	if err != nil {
		_ = g.AbortWithError(http.StatusInternalServerError, fmt.Errorf("read url response: %w", err))
		return
	}

	// return response
	g.Data(res.StatusCode(), contentType, rb)
}
