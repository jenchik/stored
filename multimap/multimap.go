package multimap

import (
	"github.com/jenchik/stored/api"
	"sync/atomic"
)

type mMap struct {
	maps      []api.StoredMap
	factory   api.StoredCopier
	instances uint64
	counter   uint64
	safe      bool
}

func newMultiMap(instances uint64, safe bool, factory api.StoredCopier) api.StoredMap {
	sm := &mMap{
		maps:      make([]api.StoredMap, instances),
		factory:   factory,
		instances: instances,
		safe:      safe,
	}
	for i := uint64(0); i < instances; i++ {
		sm.maps[i] = sm.factory.Copy()
	}
	return sm
}

// If instances < 0, then setted mode is unsafe
func New(instances int, factory api.StoredCopier) api.StoredMap {
	var unsafe bool
	if instances < 0 {
		unsafe = true
		instances *= -1
	}
	if instances < 2 {
		instances = 2
	} else if instances > 1000 {
		instances = 1000
	}
	return newMultiMap(uint64(instances), !unsafe, factory)
}

func (sm *mMap) Atomic(f api.AtomicFunc) {
	if !sm.safe {
		ptr := atomic.LoadUint64(&sm.counter) % sm.instances
		sm.maps[ptr].Atomic(f)
		return
	}
	m := newMapper(sm)
	go f(m)
}

func (sm *mMap) AtomicWait(f api.AtomicFunc) {
	if !sm.safe {
		ptr := atomic.AddUint64(&sm.counter, 1) % sm.instances
		sm.maps[ptr].AtomicWait(f)
		return
	}
	m := newMapper(sm)
	f(m)
}

func (sm *mMap) Find(key string) (value interface{}, found bool) {
	ptr := atomic.AddUint64(&sm.counter, 1) % sm.instances
	if v, f := sm.maps[ptr].Find(key); f || !sm.safe {
		return v, f
	}
	for i, _ := range sm.maps {
		if ptr == uint64(i) {
			continue
		}
		if v, f := sm.maps[i].Find(key); f {
			return v, f
		}
	}
	return nil, false
}

func (sm *mMap) Insert(key string, value interface{}) {
	if !sm.safe {
		ptr := atomic.AddUint64(&sm.counter, 1) % sm.instances
		sm.maps[ptr].Insert(key, value)
		return
	}
	go func() {
		for i, _ := range sm.maps {
			sm.maps[i].Insert(key, value)
		}
	}()
}

func (sm *mMap) Delete(key string) {
	if !sm.safe {
		ptr := atomic.AddUint64(&sm.counter, 1) % sm.instances
		sm.maps[ptr].Delete(key)
		return
	}
	go func() {
		for i, _ := range sm.maps {
			sm.maps[i].Delete(key)
		}
	}()
}

func (sm *mMap) Len() int {
	ptr := atomic.AddUint64(&sm.counter, 1) % sm.instances
	return sm.maps[ptr].Len()
}

func (sm *mMap) Update(key string, updater api.UpdateFunc) {
	ptr := atomic.AddUint64(&sm.counter, 1) % sm.instances
	if !sm.safe {
		sm.maps[ptr].Update(key, updater)
		return
	}
	go func() {
		for i, _ := range sm.maps {
			sm.maps[i].Delete(key)
		}
		sm.maps[ptr].Update(key, updater)
	}()
	//sm.Insert(key, updater(nil, false))
}

func (sm *mMap) each(child api.StoredMap, c chan<- api.Mapper) {
	var index int
	child.Each(func(m api.Mapper) {
		c <- m
		index++
		if index == m.Len() {
			c <- nil
		}
	})
}

// If unsafe mode then will be received part datas
func (sm *mMap) Each(f api.ForeachFunc) {
	if !sm.safe {
		ptr := atomic.AddUint64(&sm.counter, 1) % sm.instances
		sm.maps[ptr].Each(f)
		return
	}
	c := make(chan api.Mapper, 1)
	go func() {
		for i, _ := range sm.maps {
			sm.each(sm.maps[i], c)
		}
	}()
	go func() {
		var completed int
		m := newMapperCache(sm)
		for cm := range c {
			if cm == nil {
				completed++
				if completed == len(sm.maps) {
					close(c)
				}
				continue
			}
			if _, found := m.cache[cm.Key()]; found {
				continue
			}
			m.cache[cm.Key()] = cm.Value()
		}
		for k, _ := range m.cache {
			m.key = k
			f(m)
			if m.stop {
				return
			}
		}
	}()
}

func (sm *mMap) Copy() api.StoredMap {
	return newMultiMap(sm.instances, sm.safe, sm.factory)
}
