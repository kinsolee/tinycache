/**
  Tinycache implements a tiny memory cache library.
*/

package tinycache

import (
	"github.com/c2h5oh/datasize"
	"sort"
	"sync"
	"time"
)

/**
Cache interface that supports expire time and maximum memory setting.
*/
type Cache interface {
	// Size is a string, support argument like: 1KB，100KB，1MB，2MB，1GB
	SetMaxMemory(size string) bool
	// Set a cache item and expires after the specific time
	Set(key string, val interface{}, expire time.Duration)
	// Get a cache item
	Get(key string) (interface{}, bool)
	// Delete a cache item
	Del(key string) bool
	// Check if a cache item is exists
	Exists(key string) bool
	// Flush all cache items
	Flush() bool
	// Calculate the number of cache keys
	Keys() int64
}

type tinyCache struct {
	sync.RWMutex

	// The max memory size
	maxMem datasize.ByteSize

	// All cached items
	items map[string]*cacheItem

	// The queue orders the items by their expire time
	queue *priorityQueue

	// A signal channel triggers the flush of the timer in process of expiration
	itemSignal chan struct{}
}

func New() *tinyCache {
	cache := &tinyCache{
		maxMem:     100 * datasize.MB, // Default to 100M
		items:      make(map[string]*cacheItem),
		queue:      newPriorityQueue(),
		itemSignal: make(chan struct{}),
	}

	go cache.processExpiration()

	return cache
}

func (c *tinyCache) processExpiration() {
	calculateWaitDuration := func(nextExpiredAt time.Time) time.Duration {
		if nextExpiredAt.IsZero() {
			return time.Hour
		}
		nextDuration := time.Until(nextExpiredAt)

		if nextDuration < 0 {
			return time.Microsecond
		}
		return nextDuration
	}

	timer := time.NewTimer(time.Hour)

	for {
		// No need to deal with the queue
		if c.queue.Len() <= 0 {
			timer.Stop()
			continue
		}
		sort.Sort(c.queue)
		nextExpiredAt := (*c.queue)[0].expireAt
		waitDuration := calculateWaitDuration(nextExpiredAt)
		timer.Reset(waitDuration)

		select {
		case <-timer.C:
			timer.Stop()
			c.evictItems()
		case <-c.itemSignal:
			timer.Stop()
		}
	}
}

func (c *tinyCache) evictItems() {
	c.Lock()
	defer c.Unlock()

	if c.queue.Len() == 0 {
		return
	}

	for _, item := range *c.queue {
		if !item.expired() || item.expire == 0 {
			break
		}
		c.queue.removeItem(item)
		delete(c.items, item.key)
	}

}

func (c *tinyCache) SetMaxMemory(size string) bool {
	c.Lock()
	defer c.Unlock()
	// Set the value back to previous if error or zero
	prevMax := c.maxMem
	err := c.maxMem.UnmarshalText([]byte(size))
	if err != nil || c.maxMem == 0 {
		c.maxMem = prevMax
		return false
	}
	return true
}

func (c *tinyCache) Set(key string, val interface{}, expire time.Duration) {
	c.Lock()
	defer c.Unlock()
	defer func() {
		c.itemSignal <- struct{}{}
	}()
	item, exists := c.items[key]
	if exists {
		item.update(val, expire)
		c.queue.updateItem(item)
		return
	}
	item = newCacheItem(key, val, expire)
	c.items[key] = item
	c.queue.pushItem(item)
}

func (c *tinyCache) Get(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()
	val, ok := c.items[key]
	return val, ok
}

func (c *tinyCache) Del(key string) bool {
	c.Lock()
	defer c.Unlock()

	defer func() {
		c.itemSignal <- struct{}{}
	}()

	item, exists := c.items[key]
	if !exists {
		return false
	}
	delete(c.items, key)
	c.queue.removeItem(item)
	return true
}

func (c *tinyCache) Exists(key string) bool {
	c.RLock()
	defer c.RUnlock()
	_, exists := c.items[key]
	return exists
}

func (c *tinyCache) Flush() bool {
	c.Lock()
	defer c.Unlock()

	c.items = make(map[string]*cacheItem)
	c.queue = newPriorityQueue()
	return true
}

func (c *tinyCache) Keys() int64 {
	c.RLock()
	defer c.RUnlock()
	return int64(len(c.items))
}
