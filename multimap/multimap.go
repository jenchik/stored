package multimap

import (
	"github.com/jenchik/stored/api"
	"sync/atomic"
)

type mMap struct {
	maps      []api.StoredMap
	factory   api.StoredCopier
	queue     chan func()
	instances uint64
	counter   uint64
	safe      bool
	mapper    *mapItem
}

// If instances < 0, then setted mode is unsafe
func New(instances int, factory api.StoredCopier) api.StoredMap {
	return NewN(instances, factory, 2000)
}

// If instances < 0, then setted mode is unsafe
func NewN(instances int, factory api.StoredCopier, queueSize int) api.StoredMap {
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
	return newMultiMap(uint64(instances), !unsafe, factory, queueSize)
}

func newMultiMap(instances uint64, safe bool, factory api.StoredCopier, queueSize int) api.StoredMap {
	sm := &mMap{
		maps:      make([]api.StoredMap, instances),
		factory:   factory,
		queue:     make(chan func(), queueSize),
		instances: instances,
		safe:      safe,
	}
	sm.mapper = newMapper(sm)
	for i := uint64(0); i < instances; i++ {
		sm.maps[i] = sm.factory.Copy()
	}
	go sm.worker()
	return sm
}

func (sm *mMap) worker() {
	for f := range sm.queue {
		f()
	}
}

func (sm *mMap) Atomic(f api.AtomicFunc) {
	if !sm.safe {
		ptr := atomic.LoadUint64(&sm.counter) % sm.instances
		sm.maps[ptr].Atomic(f)
		return
	}
	sm.queue <- func() {
		sm.mapper.reset()
		f(sm.mapper)
	}
}

func (sm *mMap) AtomicWait(f api.AtomicFunc) {
	if !sm.safe {
		ptr := atomic.AddUint64(&sm.counter, 1) % sm.instances
		sm.maps[ptr].AtomicWait(f)
		return
	}
	reply := make(chan struct{})
	sm.queue <- func() {
		sm.mapper.reset()
		f(sm.mapper)
		reply <- struct{}{}
	}
	<-reply
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
	sm.queue <- func() {
		for i, _ := range sm.maps {
			sm.maps[i].Insert(key, value)
		}
	}
}

func (sm *mMap) Delete(key string) {
	if !sm.safe {
		ptr := atomic.AddUint64(&sm.counter, 1) % sm.instances
		sm.maps[ptr].Delete(key)
		return
	}
	sm.queue <- func() {
		for i, _ := range sm.maps {
			sm.maps[i].Delete(key)
		}
	}
}

func (sm *mMap) Len() int {
	ptr := atomic.AddUint64(&sm.counter, 1) % sm.instances
	return sm.maps[ptr].Len()
}

func (sm *mMap) Update(key string, updater api.UpdateFunc) {
	if !sm.safe {
		ptr := atomic.AddUint64(&sm.counter, 1) % sm.instances
		sm.maps[ptr].Update(key, updater)
		return
	}
	sm.queue <- func() {
		value := updater(sm.Find(key))
		for i, _ := range sm.maps {
			sm.maps[i].Insert(key, value)
		}
	}
}

func (sm *mMap) each(child api.StoredMap, c chan<- string, done chan<- struct{}) {
	var index int
	child.Each(func(m api.Mapper) {
		c <- m.Key()
		index++
		if index == m.Len() {
			done <- struct{}{}
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
	c := make(chan string, len(sm.maps))
	done := make(chan struct{}, len(sm.maps))
	sm.queue <- func() {
		func() {
			for i, _ := range sm.maps {
				go sm.each(sm.maps[i], c, done)
			}
		}()
		var completed int
		var key string
		m := newMapperCache(sm)
	LOOP:
		for {
			select {
			case key = <-c:
				if _, found := m.cache[key]; !found {
					m.cache[key] = struct{}{}
				}
			case <-done:
				completed++
				if completed == len(sm.maps) {
					break LOOP
				}
			}
		}
		for k, _ := range m.cache {
			m.key = k
			f(m)
			if m.stop {
				return
			}
		}
	}
}

func (sm *mMap) Copy() api.StoredMap {
	return newMultiMap(sm.instances, sm.safe, sm.factory, cap(sm.queue))
}
