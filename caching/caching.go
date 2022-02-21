package caching

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var SystemCache = cache.New(6*time.Hour, 30*time.Minute)
