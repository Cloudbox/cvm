package web

import (
	"time"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/Cloudbox/cvm/logger"
)

type Client struct {
	cacheHours      int
	maxResponseSize int64

	store *persist.MemoryStore
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

		store: persist.NewMemoryStore(time.Second),
		log:   logger.New(""),
	}
}

func (c *Client) SetHandlers(r *gin.Engine) {
	// core
	r.GET("/version", cache.CacheByRequestURI(c.store, time.Duration(c.cacheHours)*time.Hour), c.Version)
}
