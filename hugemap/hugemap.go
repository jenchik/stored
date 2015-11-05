package hugemap

import (
	"github.com/jenchik/stored/api"
	"github.com/jenchik/stored/iterator"
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

type safeMap struct {
	store map[string]interface{}
	c     chan commandData
	lock  *sync.RWMutex
}

func New() api.StoredMap {
	sm := safeMap{
		store: make(map[string]interface{}),
		c:     make(chan commandData),
		lock:  new(sync.RWMutex),
	}
	go sm.run()
	return &sm
}

func safeAtomic(m *mapItem, f api.AtomicFunc) {
	defer func() {
		m.Stop()
		m.tryUnlock()
	}()
	f(m)
}

func (sm *safeMap) run() {
	mapper := newMapper(sm)
	for command := range sm.c {
		switch command.action {
		case atomic:
			mapper.reset()
			safeAtomic(mapper, command.fatomic)
		case atomicWait:
			mapper.reset()
			safeAtomic(mapper, command.fatomic)
			command.result <- struct{}{}
		case insert:
			sm.lock.Lock()
			sm.store[command.key] = command.value
			sm.lock.Unlock()
		case remove:
			sm.lock.Lock()
			delete(sm.store, command.key)
			sm.lock.Unlock()
		case each:
			// deprecated
			mapper.reset()
			mapper.it = iterator.NewDummy()
			mapper.RLock()
			for mapper.key = range sm.store {
				command.foreach(mapper)
				if mapper.done {
					break
				}
			}
			mapper.tryUnlock()
		case update:
			value, found := sm.store[command.key]
			sm.lock.Lock()
			sm.store[command.key] = command.updater(value, found)
			sm.lock.Unlock()
		}
	}
}

func (sm *safeMap) Atomic(f api.AtomicFunc) {
	if f == nil {
		return
	}
	sm.c <- commandData{action: atomic, fatomic: f}
}

func (sm *safeMap) AtomicWait(f api.AtomicFunc) {
	if f == nil {
		return
	}
	reply := make(chan interface{})
	sm.c <- commandData{action: atomicWait, fatomic: f, result: reply}
	<-reply
}

func (sm *safeMap) Find(key string) (value interface{}, found bool) {
	sm.lock.RLock()
	defer sm.lock.RUnlock()
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
	sm.lock.RLock()
	defer sm.lock.RUnlock()
	return len(sm.store)
}

func (sm *safeMap) Update(key string, updater api.UpdateFunc) {
	if updater == nil {
		return
	}
	sm.c <- commandData{action: update, key: key, updater: updater}
}

func (sm *safeMap) Each(f api.ForeachFunc) {
	if f == nil {
		return
	}
	sm.c <- commandData{action: each, foreach: f}
}

func (sm *safeMap) Copy() api.StoredMap {
	return New()
}
