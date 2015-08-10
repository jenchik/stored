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
}

func newMultiMap(instances uint64, factory api.StoredCopier) api.StoredMap {
	sm := &mMap{
		maps:      make([]api.StoredMap, instances),
		factory:   factory,
		instances: instances,
	}
	for i:= uint64(0); i < instances; i++ {
		sm.maps[i] = sm.factory.Copy()
	}
    return sm
}

func New(instances int, factory api.StoredCopier) api.StoredMap {
	if instances < 2 {
		instances = 2
	} else if instances > 1000 {
		instances = 1000
	}
	return newMultiMap(uint64(instances), factory)
}

func (sm *mMap) Atomic(f api.AtomicFunc) {
	ptr := atomic.LoadUint64(&sm.counter) % sm.instances
	sm.maps[ptr].Atomic(f)
}

func (sm *mMap) AtomicWait(f api.AtomicFunc) {
	ptr := atomic.AddUint64(&sm.counter, 1) % sm.instances
	sm.maps[ptr].AtomicWait(f)
}

func (sm *mMap) Find(key string) (value interface{}, found bool) {
	ptr := atomic.AddUint64(&sm.counter, 1) % sm.instances
	return sm.maps[ptr].Find(key)
}

func (sm *mMap) Insert(key string, value interface{}) {
	ptr := atomic.AddUint64(&sm.counter, 1) % sm.instances
	sm.maps[ptr].Insert(key, value)
}

func (sm *mMap) Delete(key string) {
	ptr := atomic.AddUint64(&sm.counter, 1) % sm.instances
	sm.maps[ptr].Delete(key)
}

func (sm *mMap) Len() int {
	ptr := atomic.AddUint64(&sm.counter, 1) % sm.instances
	return sm.maps[ptr].Len()
}

func (sm *mMap) Update(key string, updater api.UpdateFunc) {
	ptr := atomic.AddUint64(&sm.counter, 1) % sm.instances
	sm.maps[ptr].Update(key, updater)
}

func (sm *mMap) Each(f api.ForeachFunc) {
	ptr := atomic.AddUint64(&sm.counter, 1) % sm.instances
	sm.maps[ptr].Each(f)
}

func (sm *mMap) Copy() api.StoredMap {
	return newMultiMap(sm.instances, sm.factory)
}
