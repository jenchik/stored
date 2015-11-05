package multimap

import (
	"github.com/jenchik/stored/api"
)

type mapItem struct {
	sm   *mMap
	done bool
	it   api.Iterator
	key  string
}

var _ api.Mapper = &mapItem{}

func newMapper(sm *mMap) *mapItem {
	return &mapItem{
		sm: sm,
	}
}

func (m *mapItem) do(f func(api.StoredMap) bool) {
	for i, _ := range m.sm.maps {
		if f(m.sm.maps[i]) {
			return
		}
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
	m.do(func(cm api.StoredMap) bool {
		value, found = cm.Find(key)
		return found
	})
	return
}

func (m *mapItem) Key() string {
	return m.key
}

func (m *mapItem) SetKey(key string) {
	m.key = key
}

func (m *mapItem) Value() interface{} {
	v, _ := m.Find(m.key)
	return v
}

func (m *mapItem) Delete() {
	m.do(func(cm api.StoredMap) bool {
		cm.Delete(m.key)
		return false
	})
}

func (m *mapItem) Update(value interface{}) {
	m.do(func(cm api.StoredMap) bool {
		cm.Insert(m.key, value)
		return false
	})
}

func (m *mapItem) Len() int {
	// TODO
	return m.sm.Len()
}

func (m *mapItem) Lock() {
	// TODO
}

func (m *mapItem) Unlock() {
	// TODO
}

func (m *mapItem) Clear() {
	// TODO
}

func (m *mapItem) Close() {
	// deprecated
}
