package hugemap

import (
	"sync"

	"github.com/jenchik/stored/api"
)

var _ api.Mapper = &mapItem{}

type mapItem struct {
	sm   *safeMap
	key  string
	stop bool
	lock int // 0 - no; 1 - write/read; >1 - read; <0 - disabled
	mux  sync.RWMutex
}

func newMapper(sm *safeMap) *mapItem {
	mi := mapItem{sm: sm}
	return &mi
}

func (m *mapItem) doWrite(f func()) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	if m.lock != 1 {
		m.sm.lock.Lock()
		defer m.sm.lock.Unlock()
	}
	f()
}

func (m *mapItem) reset() {
	m.key = ""
	m.lock = 0
	m.stop = false
}

func (m *mapItem) Find(key string) (value interface{}, found bool) {
	m.sm.lock.RLock()
	value, found = m.sm.store[key]
	m.sm.lock.RUnlock()
	return
}

func (m *mapItem) Key() string {
	return m.key
}

func (m *mapItem) SetKey(key string) {
	m.key = key
}

func (m *mapItem) Value() (v interface{}) {
	v, _ = m.Find(m.key)
	return
}

func (m *mapItem) Delete() {
	m.doWrite(func() {
		delete(m.sm.store, m.key)
	})
}

func (m *mapItem) Update(value interface{}) {
	m.doWrite(func() {
		m.sm.store[m.key] = value
	})
}

func (m *mapItem) Len() (n int) {
	m.sm.lock.RLock()
	n = len(m.sm.store)
	m.sm.lock.RUnlock()
	return
}

func (m *mapItem) Lock() {
	m.mux.Lock()
	defer m.mux.Unlock()
	if m.lock == 1 {
		return
	}
	m.sm.lock.Lock()
	m.lock = 1
}

func (m *mapItem) Unlock() {
	m.mux.Lock()
	defer m.mux.Unlock()
	if m.lock != 1 {
		return
	}
	m.sm.lock.Unlock()
	m.lock = 0
}

func (m *mapItem) Stop() {
	m.stop = true
}

func (m *mapItem) Clear() {
	m.doWrite(func() {
		m.sm.store = make(map[string]interface{})
	})
}

func (m *mapItem) Close() {
	// TODO
	close(m.sm.c)
}
