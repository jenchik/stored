package simplemap

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/jenchik/stored/api"
	"github.com/jenchik/stored/test"
)

func TestInsertMethods(t *testing.T) {
	sm := New()
	err := test.InserterBasic(sm, "Insert")
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

func TestAtomicMethods(t *testing.T) {
	sm := New()
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
			case c <- &test.Item{K: k, V: v}:
			}
			sm.Atomic(func(mp api.Mapper) {
				t := <-c
				if _, found := mp.Find(t.K); found {
					stop <- fmt.Errorf("Key '%s' is duplicated.", t.K)
					return
				}
				mp.SetKey(t.K)
				mp.Update(t.V)
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
	for {
		if sm.Len() >= test.CntWorks*test.CntItems {
			break
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
	sm := New()
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

func TestUpdateMethods(t *testing.T) {
	sm := New()
	updater := func(args ...interface{}) error {
		m, ok := args[0].(map[string]string)
		if !ok {
			return fmt.Errorf("Get error type 'Map'")
		}
		stop := make(chan error, 1)
		c := make(chan test.Item, len(m))
		for k, v := range m {
			select {
			case err := <-stop:
				return err
			case c <- test.Item{K: k, V: v}:
			}
			sm.Update(k, func(value interface{}, found bool) interface{} {
				t := <-c
				if found {
					stop <- fmt.Errorf("Key '%s' is duplicated.", t.K)
				}
				return t.V
			})
		}
		return nil
	}
	err := test.DoPools(func(w *test.Worker) {
		for i := range test.Data {
			w.Add(updater, test.Data[i])
		}
	}, len(test.Data), "Update")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if sm.Len() != test.CntWorks*test.CntItems {
		t.Fatal("Not equal.")
	}
	if err := test.FinderBasic(sm); err != nil {
		t.Fatalf(err.Error())
	}

	rndKeys := test.Data[rand.Intn(len(test.Data)-1)]
	newData := make(map[string]string, len(rndKeys))
	stop := make(chan error, 1)
	c := make(chan test.Item, len(rndKeys))
	for k, _ := range rndKeys {
		tt := test.Item{K: k, V: test.RandString(test.SizeItem)}
		select {
		case err := <-stop:
			t.Fatalf(err.Error())
		case c <- tt:
		}
		newData[k] = tt.V
		sm.Update(k, func(value interface{}, found bool) interface{} {
			t := <-c
			if !found {
				stop <- fmt.Errorf("Key '%s' not found.", t.K)
			}
			return t.V
		})
	}
	for k, v := range newData {
		if val, found := sm.Find(k); !found || val.(string) != v {
			t.Fatalf("Cannot found!")
		}
	}
}

func TestEachMethods(t *testing.T) {
	sm := New()
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
		if index == mp.Len() {
			stop <- nil
		}
	})
	err = <-stop
	close(stop)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestDeleteMethods(t *testing.T) {
	sm := New()
	err := test.InserterBasic(sm, "Delete")
	if err != nil {
		t.Fatalf(err.Error())
	}

	rndKeys := test.Data[rand.Intn(len(test.Data)-1)]
	for k, _ := range rndKeys {
		sm.Delete(k)
	}
	for k, _ := range rndKeys {
		if _, found := sm.Find(k); found {
			t.Fatalf("Key '%s' not deleted!", k)
		}
	}
}
