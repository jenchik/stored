package multimap

import (
	"github.com/jenchik/stored/api"
)

type mapItemCache struct {
	sm    *mMap
	key   string
	stop  bool
	cache map[string]struct{}
}

var _ api.Mapper = &mapItemCache{}

func newMapperCache(sm *mMap) *mapItemCache {
	return &mapItemCache{
		sm:    sm,
		cache: make(map[string]struct{}, sm.Len()),
	}
}

func (m *mapItemCache) find(key string) (value interface{}, found bool) {
	value, found = m.sm.Find(key)
	return
}

func (m *mapItemCache) do(f func(api.StoredMap) bool) {
	for i, _ := range m.sm.maps {
		f(m.sm.maps[i])
	}
}

func (m *mapItemCache) Next() bool {
	return false
}

func (m *mapItemCache) Stop() {
	m.stop = true
}

func (m *mapItemCache) Find(key string) (value interface{}, found bool) {
	return m.find(key)
}

func (m *mapItemCache) Key() string {
	return m.key
}

func (m *mapItemCache) SetKey(key string) {
	m.key = key
}

func (m *mapItemCache) Value() interface{} {
	v, _ := m.find(m.key)
	return v
}

func (m *mapItemCache) Delete() {
	delete(m.cache, m.key)
	m.do(func(cm api.StoredMap) bool {
		cm.Delete(m.key)
		return false
	})
}

func (m *mapItemCache) Update(value interface{}) {
	m.do(func(cm api.StoredMap) bool {
		cm.Insert(m.key, value)
		return false
	})
}

func (m *mapItemCache) Len() int {
	return len(m.cache)
}

func (m *mapItemCache) Lock() {
	// TODO
}

func (m *mapItemCache) Unlock() {
	// TODO
}

func (m *mapItemCache) Clear() {
	// TODO
}

func (m *mapItemCache) Close() {
	// deprecated
}
