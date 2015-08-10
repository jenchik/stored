package stored

import (
	"github.com/jenchik/stored/api"
	"sync/atomic"
	"time"
)

type WBCache struct {
	*Cache
}

func NewWBCache(opt Options, sm api.StoredMap) GCCacher {
	hc := &WBCache{NewCache(opt, sm)}
	return hc
}

func (c *WBCache) GetItem(key string, fback api.GetterFunc) (interface{}, error) {
	if value, found := c.sm.Find(key); found {
		valItem := value.(*item)
		if time.Now().Unix() <= valItem.ttl {
			atomic.AddUint64(&c.cntHit, 1)
			return valItem.val, nil
		}
	}
	return c.Cache.GetItem(key, fback)
}
