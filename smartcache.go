package smartcache

import (
	"context"
	"sync"
	"time"

	"github.com/Songmu/flextime"
	"golang.org/x/sync/singleflight"
)

// Cache for cache
type Cache struct {
	expire     time.Duration
	softExpire time.Duration
	generator  func(context.Context) (interface{}, error)

	group singleflight.Group

	mu             sync.RWMutex
	value          interface{}
	nextSoftExpire time.Time
	nextExpire     time.Time
}

// New returns new Cache
func New(expire, softExpire time.Duration, gen func(context.Context) (interface{}, error)) *Cache {
	return &Cache{
		expire:     expire,
		softExpire: softExpire,
		generator:  gen,
	}
}

func (ca *Cache) renew(ctx context.Context) (interface{}, error) {
	val, err, _ := ca.group.Do("renew", func() (interface{}, error) {
		now := flextime.Now()
		if now.Before(ca.nextExpire) && (ca.nextSoftExpire.IsZero() || now.Before(ca.nextSoftExpire)) {
			return ca.value, nil
		}
		val, err := ca.generator(ctx)
		if err == nil {
			ca.mu.Lock()
			ca.value = val
			if ca.softExpire > 0 {
				ca.nextSoftExpire = now.Add(ca.softExpire)
			}
			ca.nextExpire = now.Add(ca.expire)
			ca.mu.Unlock()
		}
		return val, err
	})
	return val, err
}

// Get the cached value
func (ca *Cache) Get(ctx context.Context) (interface{}, error) {
	now := flextime.Now()
	ca.mu.RLock()
	currVal := ca.value
	softExpire := ca.nextSoftExpire
	expire := ca.nextExpire
	ca.mu.RUnlock()

	if now.After(expire) {
		return ca.renew(ctx)
	}
	if !softExpire.IsZero() && now.After(softExpire) {
		go ca.renew(ctx)
	}
	return currVal, nil
}
