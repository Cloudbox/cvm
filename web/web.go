package web

import (
	"github.com/Cloudbox/cvm/logger"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"golang.org/x/sync/singleflight"
	"time"
)

type Client struct {
	cacheHours      int
	maxResponseSize int64

	store *persistence.InMemoryStore
	sfg   *singleflight.Group
	log   zerolog.Logger
}

type Config struct {
	CacheHours      int   `yaml:"cache_hours"`
	MaxResponseSize int64 `yaml:"max_response_size"`
}

func New(c *Config) *Client {
	if c.CacheHours == 0 {
		c.CacheHours = 12
	}

	if c.MaxResponseSize == 0 {
		c.MaxResponseSize = 5
	}

	return &Client{
		cacheHours:      c.CacheHours,
		maxResponseSize: c.MaxResponseSize,

		store: persistence.NewInMemoryStore(time.Second),
		sfg:   &singleflight.Group{},
		log:   logger.New(""),
	}
}

func (c *Client) SetHandlers(r *gin.Engine) {
	// core
	r.GET("/version", cache.CachePage(c.store, time.Duration(c.cacheHours)*time.Hour, c.Version))
}
