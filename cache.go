package stored

import (
	"github.com/jenchik/stored/api"
	"sync"
	"sync/atomic"
	"time"
)

type item struct {
	val interface{}
	ttl int64
}

type Cache struct {
	sm        api.StoredMap
	maxItems  int
	ttl       int64
	cntGcCall uint64 // счётчик вызовов мусорщика
	cntGcHit  uint64 // счётчик успешной чистки каждого элемента
	cntHit    uint64 // общий счётчик общих обращений
	cntFail   uint64 // общий счётчик "тухлых" элементов
	cntInsert uint64 // общий счётчик вставок
	cntMissed uint64 // общий счётчик промашек
	gcTimeout time.Duration
	gcm       sync.RWMutex
}

// each instance may be creates two streams
func NewCache(opt Options, sm api.StoredMap) *Cache {
	c := &Cache{
		sm:        sm,
		maxItems:  saneMaxItems(opt.MaxItems),
		ttl:       int64(saneTTL(opt.TTL)),
		gcTimeout: time.Duration(saneTimeoutGC(opt.TimeoutGC)),
	}
	if !opt.DisabledGC {
		go gc(c)
	}
	return c
}

func (c *Cache) GetItem(key string, fback api.GetterFunc) (interface{}, error) {
	var valItem *item
	var err error
	c.sm.AtomicWait(func(m api.Mapper) {
		atomic.AddUint64(&c.cntHit, 1)
		m.SetKey(key)
		if value := m.Value(); value != nil {
			valItem = value.(*item)
			if time.Now().Unix() <= valItem.ttl {
				return
			}
			c.cntFail++
			valItem = nil
		}
		c.cntMissed++
		value, e := fback(key)
		if e != nil {
			err = e
			return
		}
		valItem = &item{val: value, ttl: time.Now().Unix() + c.ttl}
		m.Update(valItem)
		c.cntInsert++
	})
	if err != nil {
		return nil, err
	}
	return valItem.val, nil
}

func (c *Cache) ResetItem(key string) {
	c.sm.Delete(key)
}

func (c *Cache) Stats() map[string]uint64 {
	var s map[string]uint64
	c.gcm.RLock()
	defer c.gcm.RUnlock()
	c.sm.AtomicWait(func(m api.Mapper) {
		s = map[string]uint64{
			"gcCall": c.cntGcCall,
			"gcHit":  c.cntGcHit,
			"fail":   c.cntFail,
			"hit":    c.cntHit,
			"insert": c.cntInsert,
			"missed": c.cntMissed,
			"len":    uint64(m.Len()),
		}
	})
	return s
}

func (c *Cache) timeoutGC() time.Duration {
	return c.gcTimeout
}

func (c *Cache) GC() {
	c.gcm.Lock()
	defer c.gcm.Unlock()
	c.cntGcCall++
	index := 0
	c.sm.Each(func(m api.Mapper) {
		if m.Len() <= c.maxItems {
			m.Stop()
			return
		}
		if index > c.maxItems {
			c.cntGcHit++
			m.Delete()
			return
		}
		v := m.Value().(item)
		if time.Now().Unix() > v.ttl {
			c.cntGcHit++
			m.Delete()
		} else {
			index++
		}
		return
	})
}
