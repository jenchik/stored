package hugemap

import (
	"github.com/jenchik/stored/api"
	"sync"
)

type commandData struct {
	action  commandAction
	key     string
	value   interface{}
	result  chan<- interface{}
	data    chan<- map[string]interface{}
	updater api.UpdateFunc
	foreach api.ForeachFunc
	fatomic api.AtomicFunc
}

type commandAction int

const (
	remove commandAction = iota
	insert
	update
	each
	atomic
	atomicWait
)

type findResult struct {
	value interface{}
	found bool
}

type mapItem struct {
	sm    *safeMap
	key   string
	value interface{}
	stop  bool
	lock  int // 0 - no; 1 - write/read; >1 - read; <0 - disabled
	mux   sync.RWMutex
}

func newMapper(sm *safeMap) *mapItem {
	mi := mapItem{sm: sm}
	return &mi
}

func (m *mapItem) doWrite(f func()) {
	if m.lock > -1 {
		m.mux.RLock()
		defer m.mux.RUnlock()
		if m.lock != 1 {
			m.sm.m.Lock()
			defer m.sm.m.Unlock()
		}
	}
	f()
}

func (m *mapItem) Find(key string) (value interface{}, found bool) {
	value, found = m.sm.store[key]
	return
}

func (m *mapItem) Key() string {
	return m.key
}

func (m *mapItem) SetKey(key string) {
	m.key = key
}

func (m *mapItem) Value() interface{} {
	return m.sm.store[m.key]
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

func (m *mapItem) Len() int {
	return len(m.sm.store)
}

func (m *mapItem) Lock() {
	if m.lock < 0 {
		return
	}
	m.mux.Lock()
	defer m.mux.Unlock()
	if m.lock == 1 {
		return
	}
	m.sm.m.Lock()
	m.lock = 1
}

func (m *mapItem) Unlock() {
	if m.lock < 0 {
		return
	}
	m.mux.Lock()
	defer m.mux.Unlock()
	if m.lock != 1 {
		return
	}
	m.sm.m.Unlock()
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
	close(m.sm.c)
}

type safeMap struct {
	store map[string]interface{}
	c     chan commandData
	m     *sync.RWMutex
}

func New() api.StoredMap {
	sm := safeMap{
		store: make(map[string]interface{}),
		c:     make(chan commandData),
		m:     new(sync.RWMutex),
	}
	go sm.run()
	return &sm
}

func safeAtomic(m *mapItem, f api.AtomicFunc) {
	defer func() {
		if m.lock == 1 {
			m.sm.m.Unlock()
		}
	}()
	f(m)
}

func (sm *safeMap) run() {
	for command := range sm.c {
		switch command.action {
		case atomic:
			if command.fatomic != nil {
				mapper := newMapper(sm)
				safeAtomic(mapper, command.fatomic)
			}
		case atomicWait:
			if command.fatomic != nil {
				mapper := newMapper(sm)
				safeAtomic(mapper, command.fatomic)
			}
			command.result <- struct{}{}
		case insert:
			sm.m.Lock()
			sm.store[command.key] = command.value
			sm.m.Unlock()
		case remove:
			sm.m.Lock()
			delete(sm.store, command.key)
			sm.m.Unlock()
		case each:
			mapper := newMapper(sm)
			mapper.lock = -1
			for key, _ := range sm.store {
				mapper.key = key
				command.foreach(mapper)
				if mapper.stop {
					break
				}
			}
		case update:
			value, found := sm.store[command.key]
			sm.m.Lock()
			sm.store[command.key] = command.updater(value, found)
			sm.m.Unlock()
		}
	}
}

func (sm *safeMap) Atomic(f api.AtomicFunc) {
	sm.c <- commandData{action: atomic, fatomic: f}
}

func (sm *safeMap) AtomicWait(f api.AtomicFunc) {
	reply := make(chan interface{})
	sm.c <- commandData{action: atomicWait, fatomic: f, result: reply}
	<-reply
}

func (sm *safeMap) Find(key string) (value interface{}, found bool) {
	sm.m.RLock()
	defer sm.m.RUnlock()
	value, found = sm.store[key]
	return
}

func (sm *safeMap) Insert(key string, value interface{}) {
	sm.c <- commandData{action: insert, key: key, value: value}
}

func (sm *safeMap) Delete(key string) {
	sm.c <- commandData{action: remove, key: key}
}

func (sm *safeMap) Len() int {
	sm.m.RLock()
	defer sm.m.RUnlock()
	return len(sm.store)
}

func (sm *safeMap) Update(key string, updater api.UpdateFunc) {
	sm.c <- commandData{action: update, key: key, updater: updater}
}

func (sm *safeMap) Each(f api.ForeachFunc) {
	sm.c <- commandData{action: each, foreach: f}
}

func (sm *safeMap) Copy() api.StoredMap {
	return New()
}
