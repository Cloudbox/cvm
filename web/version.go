package web

import (
	"fmt"
	"github.com/Cloudbox/cvm/build"
	"github.com/gin-gonic/gin"
	"github.com/lucperkins/rek"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type versionRequest struct {
	URL string `form:"url" binding:"required"`
}

type versionResponse struct {
	StatusCode  int
	ContentType string
	Body        []byte
}

func (c *Client) Version(g *gin.Context) {
	// parse query
	b := new(versionRequest)
	if err := g.ShouldBindQuery(b); err != nil {
		g.AbortWithError(http.StatusBadRequest, fmt.Errorf("bind query: %w", err))
		return
	}

	// retrieve data
	v, err, _ := c.sfg.Do(b.URL, func() (interface{}, error) {
		// create request
		res, err := rek.Get(b.URL, rek.Timeout(15*time.Second), rek.UserAgent(build.UserAgent))
		if err != nil {
			return nil, fmt.Errorf("request url: %w", err)
		}
		defer res.Body().Close()

		// validate response
		if res.StatusCode() != http.StatusOK {
			return nil, fmt.Errorf("validate url response: %s", res.Status())
		}

		contentType := strings.ToLower(res.Headers()["Content-Type"])
		if !strings.Contains(contentType, "xml") && !strings.Contains(contentType, "json") {
			return nil, fmt.Errorf("validate url response content-type: %v", contentType)
		}

		// read response
		rb, err := ioutil.ReadAll(http.MaxBytesReader(nil, res.Body(), c.maxResponseSize<<20+1))
		if err != nil {
			return nil, fmt.Errorf("read url response: %w", err)
		}

		return &versionResponse{
			StatusCode:  res.StatusCode(),
			ContentType: contentType,
			Body:        rb,
		}, nil
	})
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// type cast result
	res, ok := v.(*versionResponse)
	if !ok {
		g.AbortWithError(http.StatusInternalServerError, fmt.Errorf("typecast url result"))
		return
	}

	// return response
	g.Data(res.StatusCode, res.ContentType, res.Body)
}
