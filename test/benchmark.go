package test

import (
	"testing"

	"github.com/jenchik/stored/api"
)

func BInsert(b *testing.B, sm api.StoredMap) {
	var k string
	l := len(UniqKey)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k = UniqKey[i%l]
		sm.Insert(k, k)
	}
}

func BAtomicUpdate(b *testing.B, sm api.StoredMap) {
	var k string
	l := len(UniqKey)
	inserter := func(key string) {
		sm.Atomic(func(m api.Mapper) {
			m.SetKey(key)
			m.Update(key)
		})
	}
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k = UniqKey[i%l]
		inserter(k)
	}
}

func BAtomicWaitUpdate(b *testing.B, sm api.StoredMap) {
	var k string
	l := len(UniqKey)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k = UniqKey[i%l]
		sm.AtomicWait(func(m api.Mapper) {
			m.SetKey(k)
			m.Update(k)
		})
	}
}

func BUpdate(b *testing.B, sm api.StoredMap) {
	var k string
	l := len(UniqKey)
	updater := func(key string) {
		sm.Update(key, func(value interface{}, found bool) interface{} {
			if found {
				_ = value.(string)
			}
			return key
		})
	}
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k = UniqKey[i%l]
		updater(k)
	}
}

func BAtomicComplex(b *testing.B, sm api.StoredMap) {
	var k string
	l := len(UniqKey)
	inserter := func(key string) {
		sm.Atomic(func(m api.Mapper) {
			if value, found := m.Find(key); found {
				_ = value.(string)
				return
			}
			m.SetKey(key)
			m.Update(key)
		})
	}
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k = UniqKey[i%l]
		inserter(k)
	}
}

func BAtomicWaitComplex(b *testing.B, sm api.StoredMap) {
	var k string
	l := len(UniqKey)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k = UniqKey[i%l]
		sm.AtomicWait(func(m api.Mapper) {
			if value, found := m.Find(k); found {
				_ = value.(string)
				return
			}
			m.SetKey(k)
			m.Update(k)
		})
	}
}

func BAtomicFind(b *testing.B, sm api.StoredMap) {
	var k string
	l := len(UniqKey)
	finder := func(key string) {
		sm.Atomic(func(m api.Mapper) {
			if value, found := m.Find(key); found {
				_ = value.(string)
			}
		})
	}
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k = UniqKey[i%l]
		finder(k)
	}
}

func BAtomicWaitFind(b *testing.B, sm api.StoredMap) {
	var k string
	l := len(UniqKey)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k = UniqKey[i%l]
		sm.AtomicWait(func(m api.Mapper) {
			if value, found := m.Find(k); found {
				_ = value.(string)
			}
		})
	}
}

func BFind(b *testing.B, sm api.StoredMap) {
	var k string
	l := len(UniqKey)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k = UniqKey[i%l]
		if value, found := sm.Find(k); found {
			_ = value.(string)
		}
	}
}

func BAtomicEachN(b *testing.B, sm api.StoredMap) {
	b.N = 1
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm.Atomic(func(mp api.Mapper) {
			var index int
			for mp.Next() {
				_ = mp.Value().(string)
				index++
				if index == CntItemsForEachN {
					mp.Stop()
				}
			}
		})
	}
}

func BAtomicEachShort(b *testing.B, sm api.StoredMap) {
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm.Atomic(func(mp api.Mapper) {
			for mp.Next() {
				_ = mp.Value().(string)
				mp.Stop()
			}
		})
	}
}

func BEachN(b *testing.B, sm api.StoredMap) {
	b.N = 1
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var index int
		sm.Each(func(m api.Mapper) {
			_ = m.Value().(string)
			index++
			if index == CntItemsForEachN {
				m.Stop()
			}
		})
	}
}

func BEachShort(b *testing.B, sm api.StoredMap) {
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm.Each(func(m api.Mapper) {
			_ = m.Value().(string)
			m.Stop()
		})
	}
}

func BDelete(b *testing.B, sm api.StoredMap) {
	var k string
	l := len(UniqKey)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k = UniqKey[i%l]
		sm.Delete(k)
	}
}

func BThreadsInsert(b *testing.B, sm api.StoredMap) {
	l := len(UniqKey)
	b.SetParallelism(CntBenchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var k string
		var i int
		for pb.Next() {
			k = UniqKey[i%l]
			sm.Insert(k, k)
			i++
		}
	})
}

func BThreadsAtomicUpdate(b *testing.B, sm api.StoredMap) {
	l := len(UniqKey)
	inserter := func(key string) {
		sm.Atomic(func(m api.Mapper) {
			m.SetKey(key)
			m.Update(key)
		})
	}
	b.SetParallelism(CntBenchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var k string
		var i int
		for pb.Next() {
			k = UniqKey[i%l]
			inserter(k)
			i++
		}
	})
}

func BThreadsAtomicWaitUpdate(b *testing.B, sm api.StoredMap) {
	l := len(UniqKey)
	b.SetParallelism(CntBenchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var k string
		var i int
		for pb.Next() {
			k = UniqKey[i%l]
			sm.AtomicWait(func(m api.Mapper) {
				m.SetKey(k)
				m.Update(k)
			})
			i++
		}
	})
}

func BThreadsUpdate(b *testing.B, sm api.StoredMap) {
	l := len(UniqKey)
	updater := func(key string) {
		sm.Update(key, func(value interface{}, found bool) interface{} {
			if found {
				_ = value.(string)
			}
			return key
		})
	}
	b.SetParallelism(CntBenchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var k string
		var i int
		for pb.Next() {
			k = UniqKey[i%l]
			updater(k)
			i++
		}
	})
}

func BThreadsAtomicComplex(b *testing.B, sm api.StoredMap) {
	l := len(UniqKey)
	inserter := func(key string) {
		sm.Atomic(func(m api.Mapper) {
			if value, found := m.Find(key); found {
				_ = value.(string)
				return
			}
			m.SetKey(key)
			m.Update(key)
		})
	}
	b.SetParallelism(CntBenchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var k string
		var i int
		for pb.Next() {
			k = UniqKey[i%l]
			inserter(k)
			i++
		}
	})
}

func BThreadsAtomicWaitComplex(b *testing.B, sm api.StoredMap) {
	l := len(UniqKey)
	b.SetParallelism(CntBenchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var k string
		var i int
		for pb.Next() {
			k = UniqKey[i%l]
			sm.AtomicWait(func(m api.Mapper) {
				if value, found := m.Find(k); found {
					_ = value.(string)
					return
				}
				m.SetKey(k)
				m.Update(k)
			})
			i++
		}
	})
}

func BThreadsAtomicFind(b *testing.B, sm api.StoredMap) {
	l := len(UniqKey)
	finder := func(key string) {
		sm.Atomic(func(m api.Mapper) {
			if value, found := m.Find(key); found {
				_ = value.(string)
			}
		})
	}
	b.SetParallelism(CntBenchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var k string
		var i int
		for pb.Next() {
			k = UniqKey[i%l]
			finder(k)
			i++
		}
	})
}

func BThreadsAtomicWaitFind(b *testing.B, sm api.StoredMap) {
	l := len(UniqKey)
	b.SetParallelism(CntBenchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var k string
		var i int
		for pb.Next() {
			k = UniqKey[i%l]
			sm.AtomicWait(func(m api.Mapper) {
				if value, found := m.Find(k); found {
					_ = value.(string)
				}
			})
			i++
		}
	})
}

func BThreadsFind(b *testing.B, sm api.StoredMap) {
	l := len(UniqKey)
	b.SetParallelism(CntBenchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var k string
		var i int
		for pb.Next() {
			k = UniqKey[i%l]
			if value, found := sm.Find(k); found {
				_ = value.(string)
			}
			i++
		}
	})
}

func BThreadsAtomicEachShort(b *testing.B, sm api.StoredMap) {
	b.SetParallelism(CntBenchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var i int
		for pb.Next() {
			sm.Atomic(func(m api.Mapper) {
				for m.Next() {
					_ = m.Value().(string)
					m.Stop()
					i++
				}
			})
		}
	})
}

func BThreadsEachShort(b *testing.B, sm api.StoredMap) {
	b.SetParallelism(CntBenchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var i int
		for pb.Next() {
			sm.Each(func(m api.Mapper) {
				_ = m.Value().(string)
				m.Stop()
			})
			i++
		}
	})
}

func BThreadsDelete(b *testing.B, sm api.StoredMap) {
	l := len(UniqKey)
	b.SetParallelism(CntBenchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var k string
		var i int
		for pb.Next() {
			k = UniqKey[i%l]
			sm.Delete(k)
			i++
		}
	})
}
