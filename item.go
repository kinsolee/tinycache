package tinycache

import (
	"sync"
	"time"
)

type cacheItem struct {
	sync.RWMutex

	key string

	data interface{}

	expireAt time.Time
	// The expire duration
	expire time.Duration

	index int // The index of the item in the heap
}

func newCacheItem(key string, data interface{}, expire time.Duration) *cacheItem {
	item := &cacheItem{
		key: key,
	}
	item.update(data, expire)
	return item
}

// Update fields when created or updated
func (ci *cacheItem) update(data interface{}, expire time.Duration) {
	ci.RLock()
	defer ci.RUnlock()
	ci.data = data
	ci.expire = expire
	ci.expireAt = time.Now().Add(expire)
}

func (ci *cacheItem) expired() bool {
	if ci.expire <= 0 {
		return false
	}
	return ci.expireAt.Before(time.Now())
}
