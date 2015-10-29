package simplemap

import (
	"sync"

	"github.com/jenchik/stored/api"
)

func safeAtomic(m *mapItem, f api.AtomicFunc) {
	defer func() {
		if m.lock == 1 {
			m.sm.lock.Unlock()
		}
	}()
	f(m)
}

type safeMap struct {
	store  map[string]interface{}
	lock   *sync.RWMutex
	atomic chan api.AtomicFunc
}

func New() api.StoredMap {
	return NewN(0, 2000)
}

func NewN(n, workers int) api.StoredMap {
	sm := &safeMap{
		store:  make(map[string]interface{}, n),
		lock:   new(sync.RWMutex),
		atomic: make(chan api.AtomicFunc, workers),
	}
	go sm.atomicWorker()
	return sm
}

func (sm *safeMap) atomicWorker() {
	mapper := newMapper(sm)
	for f := range sm.atomic {
		mapper.reset()
		safeAtomic(mapper, f)
	}
}

func (sm *safeMap) Delete(key string) {
	sm.lock.Lock()
	delete(sm.store, key)
	sm.lock.Unlock()
}

func (sm *safeMap) Find(key string) (value interface{}, found bool) {
	sm.lock.RLock()
	value, found = sm.store[key]
	sm.lock.RUnlock()
	return
}

func (sm *safeMap) Insert(key string, value interface{}) {
	sm.lock.Lock()
	sm.store[key] = value
	sm.lock.Unlock()
}

func (sm *safeMap) Atomic(f api.AtomicFunc) {
	sm.atomic <- f
}

func (sm *safeMap) AtomicWait(f api.AtomicFunc) {
	mapper := newMapper(sm)
	safeAtomic(mapper, f)
}

func (sm *safeMap) Len() (n int) {
	sm.lock.RLock()
	n = len(sm.store)
	sm.lock.RUnlock()
	return
}

func (sm *safeMap) Each(f api.ForeachFunc) {
	mapper := newMapperEach(sm)
	sm.lock.Lock()
	defer sm.lock.Unlock()
	for key, _ := range sm.store {
		mapper.key = key
		f(mapper)
		if mapper.stop {
			break
		}
	}
}

func (sm *safeMap) Update(key string, f api.UpdateFunc) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	value, found := sm.store[key]
	sm.store[key] = f(value, found)
}

func (sm *safeMap) Copy() api.StoredMap {
	return NewN(len(sm.store), cap(sm.atomic))
}
