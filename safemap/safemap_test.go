package safemap

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/jenchik/stored/api"
	"github.com/jenchik/stored/test"
)

var smForBenchmark api.StoredMap

func init() {
	smForBenchmark = New()
	err := test.InserterBasic(smForBenchmark, "Benchmark")
	if err != nil {
		panic(err.Error())
	}
}

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
		c := make(chan test.Item, len(m))
		for k, v := range m {
			select {
			case err := <-stop:
				return err
			case c <- test.Item{K: k, V: v}:
			}
			sm.Atomic(func(mp api.Mapper) {
				t := <-c
				if _, found := mp.Find(k); found {
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

func BenchmarkInsert(b *testing.B) {
	var k string
	sm := New()
	l := len(test.UniqKey)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k = test.UniqKey[i%l]
		sm.Insert(k, i)
	}
}

func BenchmarkAtomicUpdate(b *testing.B) {
	var k string
	var index int
	sm := New()
	l := len(test.UniqKey)
	inserter := func(key string) {
		sm.Atomic(func(m api.Mapper) {
			index++
			m.SetKey(key)
			m.Update(index)
		})
	}
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k = test.UniqKey[i%l]
		inserter(k)
	}
}

func BenchmarkAtomicWaitUpdate(b *testing.B) {
	var k string
	var index int
	sm := New()
	l := len(test.UniqKey)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k = test.UniqKey[i%l]
		sm.AtomicWait(func(m api.Mapper) {
			index++
			m.SetKey(k)
			m.Update(index)
		})
	}
}

func BenchmarkUpdate(b *testing.B) {
	var k string
	sm := New()
	l := len(test.UniqKey)
	updater := func(key string) {
		sm.Update(key, func(value interface{}, found bool) interface{} {
			return key
		})
	}
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k = test.UniqKey[i%l]
		updater(k)
	}
}

func BenchmarkAtomicComplex(b *testing.B) {
	var k string
	var index, ret int
	sm := New()
	l := len(test.UniqKey)
	inserter := func(key string) {
		sm.Atomic(func(m api.Mapper) {
			if value, found := m.Find(key); found {
				ret = value.(int)
				return
			}
			index++
			ret = index
			m.SetKey(key)
			m.Update(ret)
		})
	}
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k = test.UniqKey[i%l]
		inserter(k)
	}
}

func BenchmarkAtomicWaitComplex(b *testing.B) {
	var k string
	var index, ret int
	sm := New()
	l := len(test.UniqKey)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k = test.UniqKey[i%l]
		sm.AtomicWait(func(m api.Mapper) {
			if value, found := m.Find(k); found {
				ret = value.(int)
				return
			}
			index++
			ret = index
			m.SetKey(k)
			m.Update(ret)
		})
	}
}

func BenchmarkAtomicFind(b *testing.B) {
	var k string
	sm := smForBenchmark
	l := len(test.UniqKey)
	finder := func(key string) {
		sm.Atomic(func(m api.Mapper) {
			m.Find(key)
		})
	}
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k = test.UniqKey[i%l]
		finder(k)
	}
}

func BenchmarkAtomicWaitFind(b *testing.B) {
	var k string
	sm := smForBenchmark
	l := len(test.UniqKey)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k = test.UniqKey[i%l]
		sm.AtomicWait(func(m api.Mapper) {
			m.Find(k)
		})
	}
}

func BenchmarkFind(b *testing.B) {
	var k string
	sm := smForBenchmark
	l := len(test.UniqKey)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k = test.UniqKey[i%l]
		sm.Find(k)
	}
}

func BenchmarkEachFullCicle(b *testing.B) {
	sm := smForBenchmark
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm.Each(func(m api.Mapper) {
			_ = m.Value()
		})
	}
}

func BenchmarkEachShort(b *testing.B) {
	sm := smForBenchmark
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm.Each(func(m api.Mapper) {
			_ = m.Value()
			m.Stop()
		})
	}
}

func BenchmarkDelete(b *testing.B) {
	var k string
	sm := smForBenchmark
	l := len(test.UniqKey)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k = test.UniqKey[i%l]
		sm.Delete(k)
	}
}
