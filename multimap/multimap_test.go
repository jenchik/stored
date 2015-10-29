package multimap

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/jenchik/stored/api"
	smap "github.com/jenchik/stored/hugemap"
	"github.com/jenchik/stored/test"
)

// WARNING! Setted 'safe' mode
const testMultiN = 4

func newTest() api.StoredMap {
	return New(testMultiN, smap.New().(api.StoredCopier))
}

func testWaitN(sm api.StoredMap, n int) error {
	var stop bool
	mm := sm.(*mMap)
	ts := time.Now().UTC().Unix()
	for {
		stop = true
		for i, _ := range mm.maps {
			if mm.maps[i].Len() != n {
				stop = false
				break
			}
		}
		if stop {
			return nil
		}
		if time.Now().UTC().Unix()-ts > 2 { // wait 2 second
			return fmt.Errorf("Time elapsed.")
		}
	}
	return fmt.Errorf("Not equal.") // never
}

func TestInsertMethods(t *testing.T) {
	sm := newTest()
	err := test.InserterBasic(sm, "Insert")
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = testWaitN(sm, test.CntWorks*test.CntItems)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if err := test.FinderBasic(sm); err != nil {
		t.Fatalf(err.Error())
	}
}

func TestAtomicMethods(t *testing.T) {
	sm := newTest()
	done := make(chan *test.Item, len(test.UniqMap))
	inserter := func(args ...interface{}) error {
		m, ok := args[0].(map[string]string)
		if !ok {
			return fmt.Errorf("Get error type 'Map'")
		}
		stop := make(chan error)
		c := make(chan *test.Item, len(m))
		for k, v := range m {
			select {
			case err := <-stop:
				return err
			case c <- &test.Item{K: k, V: v, Done: done}:
			}
			sm.Atomic(func(mp api.Mapper) {
				t := <-c
				if _, found := mp.Find(t.K); found {
					stop <- fmt.Errorf("Key '%s' is duplicated.", t.K)
					return
				}
				mp.SetKey(t.K)
				mp.Update(t.V)
				t.Done <- t
			})
		}
		return nil
	}
	err := test.DoPools(func(w *test.Worker) {
		for i := range test.Data {
			w.Add(inserter, test.Data[i])
		}
	}, len(test.Data), "Atomic")
	if err != nil {
		t.Fatalf(err.Error())
	}
	var index int
	for _ = range done {
		index++
		if index >= len(test.UniqMap) {
			close(done)
		}
	}
	if sm.Len() != test.CntWorks*test.CntItems {
		t.Fatal("Not equal.")
	}
	if err := test.FinderBasic(sm); err != nil {
		t.Fatalf(err.Error())
	}
}

func TestAtomicWaitMethods(t *testing.T) {
	sm := newTest()
	inserter := func(args ...interface{}) error {
		m, ok := args[0].(map[string]string)
		if !ok {
			return fmt.Errorf("Get error type 'Map'")
		}
		var err error
		for k, v := range m {
			sm.AtomicWait(func(mp api.Mapper) {
				if _, found := mp.Find(k); found {
					err = fmt.Errorf("Key '%s' is duplicated.", k)
					return
				}
				mp.SetKey(k)
				mp.Update(v)
			})
			if err != nil {
				return err
			}
		}
		return nil
	}
	err := test.DoPools(func(w *test.Worker) {
		for i := range test.Data {
			w.Add(inserter, test.Data[i])
		}
	}, len(test.Data), "AtomicWait")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if sm.Len() != test.CntWorks*test.CntItems {
		t.Fatal("Not equal.")
	}
	if err := test.FinderBasic(sm); err != nil {
		t.Fatalf(err.Error())
	}
}

func update(sm api.StoredMap, t *test.Item) {
	sm.Update(t.K, func(value interface{}, found bool) interface{} {
		t.Done <- t
		return t.V
	})
}

func TestUpdateMethods(t *testing.T) {
	sm := newTest()
	err := test.InserterBasic(sm, "Update")
	if err != nil {
		t.Fatalf(err.Error())
	}

	rndKeys := test.Data[rand.Intn(len(test.Data)-1)]
	newData := make(map[string]string, len(rndKeys))
	done := make(chan *test.Item, len(rndKeys))
	for k, _ := range rndKeys {
		tt := &test.Item{
			K:    k,
			V:    test.RandString(test.SizeItem),
			Done: done,
		}
		update(sm, tt)
	}
	for tt := range done {
		newData[tt.K] = tt.V
		if len(newData) >= len(rndKeys) {
			break
		}
	}
	for k, v := range newData {
		if val, found := sm.Find(k); !found || val.(string) != v {
			t.Fatalf("Cannot found key '%s'", k)
		}
	}
}

func TestEachMethods(t *testing.T) {
	sm := newTest()
	err := test.InserterBasic(sm, "Each")
	if err != nil {
		t.Fatalf(err.Error())
	}

	stop := make(chan error, 1)
	var index int
	sm.Each(func(mp api.Mapper) {
		if mp.Len() != len(test.UniqMap) {
			stop <- fmt.Errorf("Not equal.")
			mp.Stop()
			return
		}
		if v, found := test.UniqMap[mp.Key()]; !found || mp.Value().(string) != v {
			stop <- fmt.Errorf("Key '%s' not found.", mp.Key())
			mp.Stop()
			return
		}
		index++
		if index >= len(test.UniqMap) {
			//if index == mp.Len() {
			stop <- nil
		}
	})
	err = <-stop
	close(stop)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestDeleteMethods(t *testing.T) {
	sm := newTest()
	err := test.InserterBasic(sm, "Delete")
	if err != nil {
		t.Fatalf(err.Error())
	}

	rndKeys := test.Data[rand.Intn(len(test.Data)-1)]
	for k, _ := range rndKeys {
		sm.Delete(k)
	}
	err = testWaitN(sm, len(test.UniqMap)-len(rndKeys))
	if err != nil {
		t.Fatalf(err.Error())
	}
	for k, _ := range rndKeys {
		if _, found := sm.Find(k); found {
			t.Fatalf("Key '%s' not deleted!", k)
		}
	}
}
