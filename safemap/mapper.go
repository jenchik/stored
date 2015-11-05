package safemap

import (
	"github.com/jenchik/stored/api"
	"github.com/jenchik/stored/iterator"
)

var _ api.Mapper = &mapItem{}

func newIterator(m *mapItem) api.Iterator {
	return iterator.New(m.store, &m.key)
}

type mapItem struct {
	sm    safeMap
	store map[string]interface{}
	done  bool
	it    api.Iterator
	key   string
}

func newMapper(sm safeMap) *mapItem {
	return &mapItem{
		sm: sm,
	}
}

func (m *mapItem) reset() {
	m.it = nil
	m.done = false
	m.key = ""
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
	value, found = m.store[key]
	return
}

func (m *mapItem) Key() string {
	return m.key
}

func (m *mapItem) SetKey(key string) {
	m.key = key
}

func (m *mapItem) Value() interface{} {
	return m.store[m.key]
}

func (m *mapItem) Delete() {
	delete(m.store, m.key)
}

func (m *mapItem) Update(value interface{}) {
	m.store[m.key] = value
}

func (m *mapItem) Len() int {
	return len(m.store)
}

func (m *mapItem) Lock() {
}

func (m *mapItem) Unlock() {
}

func (m *mapItem) RLock() {
}

func (m *mapItem) RUnlock() {
}

func (m *mapItem) Clear() {
	m.store = make(map[string]interface{})
}

func (m *mapItem) Close() {
	// deprecated
}
