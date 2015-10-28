package multimap

import (
	//"fmt"
	//"math/rand"
	"testing"
	//"time"

	"github.com/jenchik/stored/api"
	"github.com/jenchik/stored/safemap"
	"github.com/jenchik/stored/test"
)

// WARNING! Setted 'unsafe' mode
const testUnsafeMultiN = -4

var smUnsafeForBenchmark api.StoredMap

func init() {
	smUnsafeForBenchmark = newTest()
	err := test.InserterBasic(smUnsafeForBenchmark, "BenchmarkUnsafe")
	if err != nil {
		panic(err.Error())
	}
	err = testWaitN(smForBenchmark, test.CntWorks*test.CntItems)
	if err != nil {
		panic(err.Error())
	}
}

func newUnsafeTest() api.StoredMap {
	return New(testUnsafeMultiN, safemap.New().(api.StoredCopier))
}

func TestUnsafeInsertMethods(t *testing.T) {
	t.Skip("Do not use! Better use Atomic...")
}

func TestUnsafeAtomicMethods(t *testing.T) {
	t.Skip("TODO")
}

func TestUnsafeAtomicWaitMethods(t *testing.T) {
	t.Skip("TODO")
}

func TestUnsafeUpdateMethods(t *testing.T) {
	// TODO
	t.Skip("TODO")
}

func TestUnsafeEachMethods(t *testing.T) {
	t.Skip("Do not use! Better use Atomic...")
}

func TestUnsafeDeleteMethods(t *testing.T) {
	t.Skip("Do not use! Better use Atomic...")
}

func BenchmarkUnsafeInsert(b *testing.B) {
	var k string
	sm := newUnsafeTest()
	l := len(test.UniqKey)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k = test.UniqKey[i%l]
		sm.Insert(k, i)
	}
}

func BenchmarkUnsafeAtomicUpdate(b *testing.B) {
	var k string
	var index int
	sm := newUnsafeTest()
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

func BenchmarkUnsafeAtomicWaitUpdate(b *testing.B) {
	var k string
	var index int
	sm := newUnsafeTest()
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

func BenchmarkUnsafeUpdate(b *testing.B) {
	var k string
	sm := newUnsafeTest()
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

func BenchmarkUnsafeAtomicComplex(b *testing.B) {
	var k string
	var index, ret int
	sm := newUnsafeTest()
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

func BenchmarkUnsafeAtomicWaitComplex(b *testing.B) {
	var k string
	var index, ret int
	sm := newUnsafeTest()
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

func BenchmarkUnsafeAtomicFind(b *testing.B) {
	var k string
	sm := smUnsafeForBenchmark
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

func BenchmarkUnsafeAtomicWaitFind(b *testing.B) {
	var k string
	sm := smUnsafeForBenchmark
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

func BenchmarkUnsafeFind(b *testing.B) {
	var k string
	sm := smUnsafeForBenchmark
	l := len(test.UniqKey)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k = test.UniqKey[i%l]
		sm.Find(k)
	}
}

func BenchmarkUnsafeEachFullCicle(b *testing.B) {
	b.Skip("TODO")
	return
	sm := smUnsafeForBenchmark
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm.Each(func(m api.Mapper) {
			_ = m.Value()
		})
	}
}

func BenchmarkUnsafeEachShort(b *testing.B) {
	b.Skip("TODO")
	return
	sm := smUnsafeForBenchmark
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

func BenchmarkUnsafeDelete(b *testing.B) {
	var k string
	sm := smUnsafeForBenchmark
	l := len(test.UniqKey)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k = test.UniqKey[i%l]
		sm.Delete(k)
	}
}
