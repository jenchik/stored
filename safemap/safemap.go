// Copyright © 2011-12 Qtrac Ltd.
//
// This program or package and any associated files are licensed under the
// Apache License, Version 2.0 (the "License"); you may not use these files
// except in compliance with the License. You can get a copy of the License
// at: http://www.apache.org/licenses/LICENSE-2.0.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package safemap

import (
	"github.com/jenchik/stored/api"
)

type safeMap chan commandData

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
	find
	insert
	length
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

func New() api.StoredMap {
	sm := make(safeMap) // type safeMap chan commandData
	go sm.run()
	return sm
}

func (sm safeMap) run() {
	store := make(map[string]interface{})
	for command := range sm {
		switch command.action {
		case atomic:
			if command.fatomic != nil {
				mapper := &mapItem{store: store}
				command.fatomic(mapper)
			}
		case atomicWait:
			if command.fatomic != nil {
				mapper := &mapItem{store: store}
				command.fatomic(mapper)
			}
			command.result <- struct{}{}
		case find:
			value, found := store[command.key]
			command.result <- findResult{value, found}
		case insert:
			store[command.key] = command.value
		case remove:
			delete(store, command.key)
		case each:
			mapper := &mapItem{store: store}
			for key, _ := range store {
				mapper.key = key
				command.foreach(mapper)
				if mapper.stop {
					break
				}
			}
		case length:
			command.result <- len(store)
		case update:
			value, found := store[command.key]
			store[command.key] = command.updater(value, found)
		}
	}
}

func (sm safeMap) Atomic(f api.AtomicFunc) {
	sm <- commandData{action: atomic, fatomic: f}
}

func (sm safeMap) AtomicWait(f api.AtomicFunc) {
	reply := make(chan interface{})
	sm <- commandData{action: atomicWait, fatomic: f, result: reply}
	<-reply
}

func (sm safeMap) Find(key string) (value interface{}, found bool) {
	reply := make(chan interface{})
	sm <- commandData{action: find, key: key, result: reply}
	result := (<-reply).(findResult)
	return result.value, result.found
}

func (sm safeMap) Insert(key string, value interface{}) {
	sm <- commandData{action: insert, key: key, value: value}
}

func (sm safeMap) Delete(key string) {
	sm <- commandData{action: remove, key: key}
}

func (sm safeMap) Len() int {
	reply := make(chan interface{})
	sm <- commandData{action: length, result: reply}
	return (<-reply).(int)
}

func (sm safeMap) Update(key string, updater api.UpdateFunc) {
	sm <- commandData{action: update, key: key, updater: updater}
}

func (sm safeMap) Each(f api.ForeachFunc) {
	sm <- commandData{action: each, foreach: f}
}

func (sm safeMap) Copy() api.StoredMap {
	return New()
}
