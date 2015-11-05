package iterator

import (
	"github.com/jenchik/stored/api"
)

type itrDummy struct {
}

func NewDummy() api.Iterator {
	return &itrDummy{}
}

func (it *itrDummy) Done() <-chan struct{} {
	return nil
}

func (it *itrDummy) Stop() {
}

func (it *itrDummy) Next() bool {
	return false
}

var _ api.Iterator = &itr{}

type itr struct {
	key    *string
	state  int
	done   chan struct{}
	next   chan string
	stop   chan struct{}
	worker func()
}

func newItr(worker func(), k *string) *itr {
	return &itr{
		key:    k,
		stop:   make(chan struct{}, 1),
		worker: worker,
	}
}

func (it *itr) run() {
	it.state = 1
	it.done = make(chan struct{}, 1)
	it.next = make(chan string, 0)
	go it.worker()
}

func (it *itr) Done() <-chan struct{} {
	return it.done
}

func (it *itr) Next() bool {
	if it.state < 1 {
		if it.state < 0 {
			return false
		}
		it.run()
	}
	for *it.key = range it.next {
		return true
	}
	it.state = -1
	return false
}

func (it *itr) Stop() {
	select {
	case it.stop <- struct{}{}:
	default:
	}
}

func New(store map[string]interface{}, ptrKey *string) api.Iterator {
	var it *itr
	worker := func() {
		var key string
	LOOP:
		for key = range store {
			select {
			case <-it.stop:
				break LOOP
			case it.next <- key:
			}
		}
		close(it.next)
		it.done <- struct{}{}
	}
	it = newItr(worker, ptrKey)
	return it
}

type itrChan struct {
	*itr
	collector <-chan string
}

func NewFromChan(s <-chan string, k *string) api.Iterator {
	itm := &itrChan{collector: s}
	itm.itr = newItr(itm.worker, k)
	return itm
}

func (it *itrChan) worker() {
	var key string
	cache := make(map[string]struct{})
LOOP:
	for key = range it.collector {
		select {
		case <-it.stop:
			break LOOP
		default:
			if _, found := cache[key]; found {
				continue
			}
			cache[key] = struct{}{}
		}
		select {
		case <-it.stop:
			break LOOP
		case it.next <- key:
		}
	}
	close(it.next)
	it.done <- struct{}{}
}
