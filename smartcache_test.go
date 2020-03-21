package smartcache_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Songmu/flextime"
	"github.com/Songmu/smartcache"
)

func TestCache_Get(t *testing.T) {
	now := flextime.Now()
	defer flextime.Set(now)()

	var counter int
	ca := smartcache.New(10*time.Second, time.Second, func(ctx context.Context) (interface{}, error) {
		// Actually time.Sleep instead of flextime to emulate long running operations
		time.Sleep(50 * time.Millisecond)
		counter++
		return counter, nil
	})

	t.Run("create new cache", func(t *testing.T) {
		v, err := ca.Get(context.Background())
		if err != nil {
			t.Errorf("error should be nil, but: %s", err)
		}
		if v.(int) != 1 {
			t.Errorf("value should be 1, but %v", v)
		}
	})

	t.Run("get from cache", func(t *testing.T) {
		v, err := ca.Get(context.Background())
		if err != nil {
			t.Errorf("error should be nil, but: %s", err)
		}
		if v.(int) != 1 {
			t.Errorf("value should be 1, but %v", v)
		}
	})

	t.Run("soft expire and renew internal value", func(t *testing.T) {
		flextime.Sleep(2 * time.Second)
		// check the concurrent cache updates won't conflict
		var (
			wg   = sync.WaitGroup{}
			para = 10
		)
		wg.Add(para)
		for i := 0; i < para; i++ {
			go func() {
				defer wg.Done()
				ca.Get(context.Background())
			}()
		}
		wg.Wait()
		v, err := ca.Get(context.Background())
		if err != nil {
			t.Errorf("error should be nil, but: %s", err)
		}
		if v.(int) != 1 {
			t.Errorf("value should be 1, but %v", v)
		}
	})

	t.Run("wait for internal value update", func(t *testing.T) {
		time.Sleep(80 * time.Millisecond) // use real time.Sleep for waiting cache update
		v, err := ca.Get(context.Background())
		if err != nil {
			t.Errorf("error should be nil, but: %s", err)
		}
		if v.(int) != 2 {
			t.Errorf("value should be 2, but %v", v)
		}
	})

	t.Run("hard expire", func(t *testing.T) {
		flextime.Sleep(11 * time.Second)
		var (
			v   interface{}
			err error
		)
		go func() {
			v, err = ca.Get(context.Background())
		}()
		time.Sleep(50 * time.Millisecond)
		// check the concurrent cache updates won't conflict
		var (
			wg   = sync.WaitGroup{}
			para = 10
		)
		wg.Add(para)
		for i := 0; i < para; i++ {
			go func() {
				defer wg.Done()
				ca.Get(context.Background())
			}()
		}
		wg.Wait()
		if err != nil {
			t.Errorf("error should be nil, but: %s", err)
		}
		if v.(int) != 3 {
			t.Errorf("value should be 3, but %v", v)
		}

		v, err = ca.Get(context.Background())
		if err != nil {
			t.Errorf("error should be nil, but: %s", err)
		}
		if v.(int) != 3 {
			t.Errorf("value should be 3, but %v", v)
		}
	})
}
