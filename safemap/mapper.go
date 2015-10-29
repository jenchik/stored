package safemap

import (
	"github.com/jenchik/stored/api"
)

var _ api.Mapper = &mapItem{}

type mapItem struct {
	sm    safeMap
	store map[string]interface{}
	key   string
	value interface{}
	stop  bool
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

func (m *mapItem) Stop() {
	m.stop = true
}

func (m *mapItem) Clear() {
	m.store = make(map[string]interface{})
}

func (m *mapItem) Close() {
	close(m.sm)
}
