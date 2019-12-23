package tinycache

import (
	"testing"
	"time"
)

func Test_cacheItem_expired(t *testing.T) {
	delay := time.Microsecond * 100
	item := newCacheItem("expire", "", delay)
	if item.expired() != false {
		t.Errorf("expire() expects true, ret: %v", item.expired())
	}
	time.Sleep(delay)
	if item.expired() != true {
		t.Errorf("expire() expects false, ret: %v", item.expired())
	}
}

func Test_cacheItem_update(t *testing.T) {
	item := newCacheItem("item", "old", time.Nanosecond)
	item.update("new", time.Hour)
	if item.data != "new" || item.expire != time.Hour || item.expireAt.Before(time.Now()) {
		t.Errorf("update cache item failed")
	}
}

