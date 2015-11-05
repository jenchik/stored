package multimap

import (
	"github.com/jenchik/stored/api"
	"github.com/jenchik/stored/iterator"
	"sync"
)

func newIterator(m *mapItem) api.Iterator {
	l := len(m.sm.maps)
	bridge := make(chan string, l)
	it := iterator.NewFromChan(bridge, &m.key)
	done := make(chan struct{}, l)
	stop := make(chan struct{}, l)
	wg := new(sync.WaitGroup)
	closer := func() {
		done <- struct{}{}
		wg.Done()
	}
	for i, _ := range m.sm.maps {
		wg.Add(1)
		go each(m.sm.maps[i], stop, bridge, closer)
	}
	go func() {
		var completed int
	LOOP:
		for {
			select {
			case <-it.Done():
				for _ = range m.sm.maps {
					stop <- struct{}{}
				}
				break LOOP
			case <-done:
				completed++
				if completed == l {
					break LOOP
				}
			}
		}
		wg.Wait()
		close(bridge)
	}()
	return it
}

func each(child api.StoredMap, stop <-chan struct{}, c chan<- string, closer func()) {
	var index int
	child.Each(func(m api.Mapper) {
		select {
		case <-stop:
			m.Stop()
			closer()
			return
		case c <- m.Key():
		}
		index++
		if index == m.Len() {
			closer()
		}
	})
}
