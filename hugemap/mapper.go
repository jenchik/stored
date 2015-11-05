package hugemap

import (
	"github.com/jenchik/stored/api"
	"github.com/jenchik/stored/iterator"
)

var _ api.Mapper = &mapItem{}

func newIterator(m *mapItem) api.Iterator {
	return iterator.New(m.sm.store, &m.key)
}

type mapItem struct {
	sm   *safeMap
	done bool
	it   api.Iterator
	key  string
	lock int // 0 - no; 1 - write/read; >1 - read; <0 - disabled
}

func newMapper(sm *safeMap) *mapItem {
	return &mapItem{
		sm: sm,
	}
}

func (m *mapItem) reset() {
	m.it = nil
	m.done = false
	m.key = ""
	m.lock = 0
}

func (m *mapItem) tryUnlock() {
	if m.lock == 1 {
		m.Unlock()
	} else if m.lock > 1 {
		m.RUnlock()
	}
}

func (m *mapItem) Next() bool {
	if m.done == true {
		return false
	}
	if m.it == nil {
		if m.Len() == 0 {
			// empty
			m.done = true
			return false
		}
		m.it = newIterator(m)
	}
	if !m.it.Next() {
		m.done = true
		return false
	}
	return true
}

func (m *mapItem) Stop() {
	if !m.done && m.it != nil {
		m.it.Stop()
	}
	m.done = true
}

func (m *mapItem) Find(key string) (value interface{}, found bool) {
	if m.lock < 1 {
		m.sm.lock.RLock()
		value, found = m.sm.store[key]
		m.sm.lock.RUnlock()
		return
	}
	value, found = m.sm.store[key]
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
	if m.lock != 1 {
		m.sm.lock.Lock()
		delete(m.sm.store, m.key)
		m.sm.lock.Unlock()
		return
	}
	delete(m.sm.store, m.key)
}

func (m *mapItem) Update(value interface{}) {
	if m.lock != 1 {
		m.sm.lock.Lock()
		m.sm.store[m.key] = value
		m.sm.lock.Unlock()
		return
	}
	m.sm.store[m.key] = value
}

func (m *mapItem) Len() (n int) {
	if m.lock < 1 {
		m.sm.lock.RLock()
		n = len(m.sm.store)
		m.sm.lock.RUnlock()
		return
	}
	n = len(m.sm.store)
	return
}

func (m *mapItem) Lock() {
	if m.lock == 1 {
		return
	}
	m.sm.lock.Lock()
	m.lock = 1
	return
}

func (m *mapItem) Unlock() {
	if m.lock != 1 {
		// ?
		return
	}
	m.sm.lock.Unlock()
	m.lock = 0
}

func (m *mapItem) RLock() {
	if m.lock > 1 {
		return
	}
	m.sm.lock.RLock()
	m.lock = 2
}

func (m *mapItem) RUnlock() {
	if m.lock < 2 {
		// ?
		return
	}
	m.sm.lock.RUnlock()
	m.lock = 0
}

func (m *mapItem) Clear() {
	if m.lock != 1 {
		m.sm.lock.Lock()
		m.sm.store = make(map[string]interface{})
		m.sm.lock.Unlock()
		return
	}
	m.sm.store = make(map[string]interface{})
}

func (m *mapItem) Close() {
	// deprecated
}
