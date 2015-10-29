package simplemap

import (
	"github.com/jenchik/stored/api"
)

var _ api.Mapper = &mapItemEach{}

type mapItemEach struct {
	sm   *safeMap
	key  string
	stop bool
}

func newMapperEach(sm *safeMap) *mapItemEach {
	return &mapItemEach{sm: sm}
}

func (m *mapItemEach) Find(key string) (value interface{}, found bool) {
	value, found = m.sm.store[key]
	return
}

func (m *mapItemEach) Key() string {
	return m.key
}

func (m *mapItemEach) SetKey(key string) {
	m.key = key
}

func (m *mapItemEach) Value() (v interface{}) {
	v, _ = m.sm.store[m.key]
	return
}

func (m *mapItemEach) Delete() {
	delete(m.sm.store, m.key)
}

func (m *mapItemEach) Update(value interface{}) {
	m.sm.store[m.key] = value
}

func (m *mapItemEach) Len() (n int) {
	return len(m.sm.store)
}

func (m *mapItemEach) Lock() {
}

func (m *mapItemEach) Unlock() {
}

func (m *mapItemEach) Stop() {
	m.stop = true
}

func (m *mapItemEach) Clear() {
	m.sm.store = make(map[string]interface{})
}

func (m *mapItemEach) Close() {
	// TODO
	close(m.sm.atomic)
}
