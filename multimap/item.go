package multimap

import (
	"github.com/jenchik/stored/api"
)

type mapItem struct {
	sm    *mMap
	key   string
	stop  bool
	child api.Mapper
	inst  int
}

var _ api.Mapper = &mapItem{}

func newMapper(sm *mMap) *mapItem {
	return &mapItem{
		sm:   sm,
		inst: -1,
	}
}

func (m *mapItem) do(f func(api.StoredMap) bool) {
	for i, _ := range m.sm.maps {
		if i == m.inst {
			continue
		}
		if f(m.sm.maps[i]) {
			return
		}
	}
}

func (m *mapItem) find(key string) (value interface{}, found bool) {
	if m.child != nil {
		if value, found = m.child.Find(key); found {
			return
		}
	}
	m.do(func(cm api.StoredMap) bool {
		value, found = cm.Find(key)
		return found
	})
	return
}

func (m *mapItem) Find(key string) (value interface{}, found bool) {
	return m.find(key)
}

func (m *mapItem) Key() string {
	return m.key
}

func (m *mapItem) SetKey(key string) {
	m.key = key
}

func (m *mapItem) Value() interface{} {
	v, _ := m.find(m.key)
	return v
}

func (m *mapItem) Delete() {
	if m.child != nil {
		m.child.Delete()
		return
	}
	m.do(func(cm api.StoredMap) bool {
		cm.Delete(m.key)
		return false
	})
}

func (m *mapItem) Update(value interface{}) {
	if m.child != nil {
		m.child.Update(value)
	}
	m.do(func(cm api.StoredMap) bool {
		cm.Insert(m.key, value)
		return false
	})
}

func (m *mapItem) Len() int {
	// TODO
	if m.child == nil {
		return m.sm.Len()
	}
	return m.child.Len()
}

func (m *mapItem) Lock() {
	// TODO
}

func (m *mapItem) Unlock() {
	// TODO
}

func (m *mapItem) Stop() {
	m.stop = true
}

func (m *mapItem) Clear() {
	// TODO
}

func (m *mapItem) Close() {
	// TODO
}
