package stored

import (
	"github.com/jenchik/stored/api"
	"sync/atomic"
)

type L2Cache struct {
	sm        api.StoredMap
	maxItems  int
	order     []string
	cntGcHit  uint64
	cntHit    uint64
	cntInsert uint64
	cntMissed uint64
}

func NewL2Cache(opt Options, sm api.StoredMap) Cacher {
	c := &L2Cache{
		sm:       sm,
		maxItems: saneMaxItems(opt.MaxItems),
	}
	c.order = make([]string, 0, c.maxItems)
	return c
}

func (c *L2Cache) GetItem(key string, fback api.GetterFunc) (interface{}, error) {
	var value interface{}
	if value, found := c.sm.Find(key); found {
		atomic.AddUint64(&c.cntHit, 1)
		return value, nil
	}
	var err error
	c.sm.AtomicWait(func(m api.Mapper) {
		atomic.AddUint64(&c.cntHit, 1)
		m.SetKey(key)
		value = m.Value()
		if value != nil {
			return
		}
		c.cntMissed++
		value, err = fback(key)
		if err != nil {
			return
		}
		if len(c.order) == c.maxItems {
			c.cntGcHit++
			m.SetKey(c.order[0])
			m.Delete()
			c.order = c.order[1:]
			m.SetKey(key)
		}
		c.order = append(c.order, key)
		m.Update(value)
		c.cntInsert++
	})
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (c *L2Cache) ResetItem(key string) {
	c.sm.Atomic(func(m api.Mapper) {
		if _, found := m.Find(key); found {
			m.SetKey(key)
			m.Delete()
			for i := 0; i < len(c.order); i++ {
				if key == c.order[i] {
					c.order = append(c.order[:i], c.order[i+1:]...)
					break
				}
			}
		}
	})
}

func (c *L2Cache) Stats() map[string]uint64 {
	var s map[string]uint64
	c.sm.AtomicWait(func(m api.Mapper) {
		s = map[string]uint64{
			"gcHit":  c.cntGcHit,
			"hit":    c.cntHit,
			"insert": c.cntInsert,
			"missed": c.cntMissed,
			"len":    uint64(m.Len()),
		}
	})
	return s
}
