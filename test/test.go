package test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/jenchik/stored/api"
	"github.com/stefantalpalaru/pool"
)

const (
	chars            = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	CntWorks         = 1000
	CntItems         = 100
	SizeItem         = 5
	CntBenchWorks    = 1000
	CntItemsForEachN = 1000
)

var (
	Data    []map[string]string
	UniqMap map[string]string
	UniqKey []string
)

type Filler func(args ...interface{}) error

type Item struct {
	K, V string
	Done chan *Item
}

type Worker struct {
	p   *pool.Pool
	err chan error
}

func (w *Worker) Add(f Filler, args ...interface{}) {
	w.p.Add(func(args ...interface{}) interface{} {
		if err := f(args...); err != nil {
			w.err <- err
		}
		return nil
	}, args...)
}

func (w *Worker) Close() {
}

func init() {
	var lock, umap sync.Mutex
	Data = make([]map[string]string, 0, CntWorks)
	UniqMap = make(map[string]string, CntWorks*CntItems)
	UniqKey = make([]string, 0, CntWorks*CntItems)
	fgen := func(sizeItem int) string {
		key := RandString(sizeItem)
		for {
			if _, found := UniqMap[key]; found {
				key = RandString(sizeItem)
				continue
			}
			break
		}
		return key
	}
	work := func(args ...interface{}) error {
		m := make(map[string]string, CntItems)
		for i := 0; i < CntItems; i++ {
			val := RandString(SizeItem)
			umap.Lock()
			k := fgen(SizeItem)
			UniqMap[k] = val
			UniqKey = append(UniqKey, k)
			umap.Unlock()
			m[k] = val
		}
		lock.Lock()
		Data = append(Data, m)
		lock.Unlock()
		return nil
	}
	DoPools(func(w *Worker) {
		for i := 0; i < CntWorks; i++ {
			w.Add(work, i)
		}
	}, CntWorks, "Init")
}

func DoPools(fillFunc func(*Worker), cntWorks int, prefix string) error {
	err := make(chan error, cntWorks)
	p := pool.New(cntWorks)
	w := &Worker{p, err}
	defer w.Close()
	fillFunc(w)
	//fmt.Printf("%s status: %#v\n", prefix, p.Status())
	p.Run()
	defer p.Stop()
	//fmt.Printf("%s status: %#v\n", prefix, p.Status())
	for i := 0; i < cntWorks; i++ {
		p.Wait()
		select {
		case e := <-err:
			return e
		default:
		}
	}
	//fmt.Printf("%s status: %#v\n", prefix, p.Status())
	return nil
}

func RandString(n int) string {
	buf := make([]byte, n)
	l := len(chars)
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < n; i++ {
		buf[i] = chars[rand.Intn(l)]
	}
	return string(buf)
}

func InserterBasic(sm api.StoredMap, prefix string) error {
	inserter := func(args ...interface{}) error {
		m, ok := args[0].(map[string]string)
		if !ok {
			return fmt.Errorf("Get error type 'Map'")
		}
		for k, v := range m {
			sm.Insert(k, v)
		}
		return nil
	}
	return DoPools(func(w *Worker) {
		for i := range Data {
			w.Add(inserter, Data[i])
		}
	}, len(Data), prefix)
}

func FinderBasic(sm api.StoredMap) error {
	finder := func(args ...interface{}) error {
		m, ok := args[0].(map[string]string)
		if !ok {
			return fmt.Errorf("Get error type 'Map'")
		}
		for k, v := range m {
			if val, found := sm.Find(k); !found || val.(string) != v {
				return fmt.Errorf("Cannot found!")
			}
		}
		return nil
	}
	return DoPools(func(w *Worker) {
		for i := range Data {
			w.Add(finder, Data[i])
		}
	}, len(Data), "Find")
}

func TAtomicEachFull(t *testing.T, sm api.StoredMap) {
	err := InserterBasic(sm, "AtomicEachFull")
	if err != nil {
		t.Fatalf(err.Error())
	}

	time.Sleep(time.Millisecond * 500)

	stop := make(chan error, 1)
	sm.Atomic(func(mp api.Mapper) {
		if mp.Len() != len(UniqMap) {
			stop <- fmt.Errorf("Not equal.")
			return
		}
		var index int
		var lastKey string
		for mp.Next() {
			index++
			lastKey = mp.Key()
			if v, found := UniqMap[lastKey]; !found || mp.Value().(string) != v {
				stop <- fmt.Errorf("Key '%s' not found.", lastKey)
				return
			}
		}
		if mp.Next() == true {
			stop <- fmt.Errorf("Next is true after stop.")
			return
		}
		if index != len(UniqMap) {
			stop <- fmt.Errorf("Invalid iteration.")
			return
		}
		if lastKey != mp.Key() {
			stop <- fmt.Errorf("Invalid last key.")
			return
		}
		stop <- nil
	})
	err = <-stop
	close(stop)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TAtomicEachN(t *testing.T, sm api.StoredMap) {
	err := InserterBasic(sm, "AtomicEachN")
	if err != nil {
		t.Fatalf(err.Error())
	}

	time.Sleep(time.Millisecond * 500)

	stop := make(chan error, 1)
	sm.Atomic(func(mp api.Mapper) {
		if mp.Len() != len(UniqMap) {
			stop <- fmt.Errorf("Not equal.")
			return
		}
		var index int
		var lastKey string
		for mp.Next() {
			index++
			if v, found := UniqMap[mp.Key()]; !found || mp.Value().(string) != v {
				stop <- fmt.Errorf("Key '%s' not found.", mp.Key())
				return
			}
			if index == CntItemsForEachN {
				mp.Stop()
				lastKey = mp.Key()
			}
		}
		if mp.Next() == true {
			stop <- fmt.Errorf("Next is true after stop.")
			return
		}
		if index != CntItemsForEachN {
			stop <- fmt.Errorf("Invalid iteration.")
			return
		}
		if lastKey != mp.Key() {
			stop <- fmt.Errorf("Invalid last key.")
			return
		}
		stop <- nil
	})
	err = <-stop
	close(stop)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TEachFull(t *testing.T, sm api.StoredMap) {
	err := InserterBasic(sm, "EachFull")
	if err != nil {
		t.Fatalf(err.Error())
	}

	time.Sleep(time.Millisecond * 500)

	stop := make(chan error, 1)
	var index int
	sm.Each(func(mp api.Mapper) {
		if mp.Len() != len(UniqMap) {
			stop <- fmt.Errorf("Not equal.")
			mp.Stop()
			return
		}
		if v, found := UniqMap[mp.Key()]; !found || mp.Value().(string) != v {
			stop <- fmt.Errorf("Key '%s' not found.", mp.Key())
			mp.Stop()
			return
		}
		index++
		if index == mp.Len() {
			stop <- nil
		}
	})
	err = <-stop
	close(stop)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if index != len(UniqKey) {
		t.Fatalf("Invalid iteration.")
	}
}

func TEachN(t *testing.T, sm api.StoredMap) {
	err := InserterBasic(sm, "EachN")
	if err != nil {
		t.Fatalf(err.Error())
	}

	time.Sleep(time.Millisecond * 500)

	stop := make(chan error, 1)
	var index int
	sm.Each(func(mp api.Mapper) {
		if mp.Len() != len(UniqMap) {
			stop <- fmt.Errorf("Not equal.")
			mp.Stop()
			return
		}
		if v, found := UniqMap[mp.Key()]; !found || mp.Value().(string) != v {
			stop <- fmt.Errorf("Key '%s' not found.", mp.Key())
			mp.Stop()
			return
		}
		index++
		if index == CntItemsForEachN {
			mp.Stop()
			stop <- nil
			return
		}
	})
	err = <-stop
	close(stop)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if index != CntItemsForEachN {
		t.Fatalf("Invalid iteration.")
	}
}
