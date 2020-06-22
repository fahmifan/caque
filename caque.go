package caque

import (
	"fmt"
	"time"

	gocache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

const defaultMaxKeys = 200
const defaultExpiration = time.Minute * 1
const defaultCleaningInterval = time.Minute * 2 // twice default expiration

// Cacher :nodoc:
type Cacher struct {
	maxKeys    int
	expiration time.Duration
	queue      *Queue
	goCache    *gocache.Cache
}

// NewCacher :nodoc:
func NewCacher(expiration, cleaningInterval time.Duration, maxKeys int) *Cacher {
	if expiration <= time.Duration(0) {
		expiration = defaultExpiration
	}

	if cleaningInterval <= time.Duration(0) {
		cleaningInterval = defaultCleaningInterval
	}

	if maxKeys <= 0 {
		maxKeys = defaultMaxKeys
	}

	cacher := &Cacher{
		expiration: expiration,
		maxKeys:    maxKeys,
		goCache:    gocache.New(expiration, cleaningInterval),
		queue:      &Queue{mapKeyToIndex: make(map[string]int)},
	}
	cacher.goCache.OnEvicted(cacher.onEvictedHandler)

	return cacher
}

// Set :nodoc:
func (c *Cacher) Set(key string, val interface{}) {
	isReachLimit := c.maxKeys == c.queue.Size()
	if isReachLimit {
		key := c.queue.Pop()
		c.goCache.Delete(key)
	}

	c.queue.Append(key)
	c.goCache.Set(key, val, c.expiration)
}

// Get :nodoc:
func (c *Cacher) Get(key string) (interface{}, bool) {
	return c.goCache.Get(key)
}

// Delete :nodoc:
func (c *Cacher) Delete(key string) {
	c.goCache.Delete(key)
	_ = c.queue.DeleteKey(key)
}

func (c *Cacher) onEvictedHandler(key string, val interface{}) {
	if ok := c.queue.DeleteKey(key); !ok {
		log.Error(fmt.Errorf("failed to delete evicted key %s ", key))
	}
}
