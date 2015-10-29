package hugemap

import (
	"testing"

	"github.com/jenchik/stored/api"
	"github.com/jenchik/stored/test"
)

var smForBenchmark api.StoredMap
var smForBenchmarkThread api.StoredMap

func init() {
	smForBenchmark = New()
	err := test.InserterBasic(smForBenchmark, "Benchmark")
	if err != nil {
		panic(err.Error())
	}
	smForBenchmarkThread = New()
	err = test.InserterBasic(smForBenchmarkThread, "BenchmarkThreads")
	if err != nil {
		panic(err.Error())
	}
}

func BenchmarkInsert(b *testing.B) {
	sm := New()
	test.BInsert(b, sm)
}

func BenchmarkAtomicUpdate(b *testing.B) {
	sm := New()
	test.BAtomicUpdate(b, sm)
}

func BenchmarkAtomicWaitUpdate(b *testing.B) {
	sm := New()
	test.BAtomicWaitUpdate(b, sm)
}

func BenchmarkUpdate(b *testing.B) {
	sm := New()
	test.BUpdate(b, sm)
}

func BenchmarkAtomicComplex(b *testing.B) {
	sm := New()
	test.BAtomicComplex(b, sm)
}

func BenchmarkAtomicWaitComplex(b *testing.B) {
	sm := New()
	test.BAtomicWaitComplex(b, sm)
}

func BenchmarkAtomicFind(b *testing.B) {
	sm := smForBenchmark
	test.BAtomicFind(b, sm)
}

func BenchmarkAtomicWaitFind(b *testing.B) {
	sm := smForBenchmark
	test.BAtomicWaitFind(b, sm)
}

func BenchmarkFind(b *testing.B) {
	sm := smForBenchmark
	test.BFind(b, sm)
}

func BenchmarkEachFullCicle(b *testing.B) {
	sm := smForBenchmark
	test.BEachFullCicle(b, sm)
}

func BenchmarkEachShort(b *testing.B) {
	sm := smForBenchmark
	test.BEachShort(b, sm)
}

func BenchmarkDelete(b *testing.B) {
	sm := smForBenchmark
	test.BDelete(b, sm)
}

func BenchmarkThreadsInsert(b *testing.B) {
	sm := New()
	test.BThreadsInsert(b, sm)
}

func BenchmarkThreadsAtomicUpdate(b *testing.B) {
	sm := New()
	test.BThreadsAtomicUpdate(b, sm)
}

func BenchmarkThreadsAtomicWaitUpdate(b *testing.B) {
	sm := New()
	test.BThreadsAtomicWaitUpdate(b, sm)
}

func BenchmarkThreadsUpdate(b *testing.B) {
	sm := New()
	test.BThreadsUpdate(b, sm)
}

func BenchmarkThreadsAtomicComplex(b *testing.B) {
	sm := New()
	test.BThreadsAtomicComplex(b, sm)
}

func BenchmarkThreadsAtomicWaitComplex(b *testing.B) {
	sm := New()
	test.BThreadsAtomicWaitComplex(b, sm)
}

func BenchmarkThreadsAtomicFind(b *testing.B) {
	sm := smForBenchmarkThread
	test.BThreadsAtomicFind(b, sm)
}

func BenchmarkThreadsAtomicWaitFind(b *testing.B) {
	sm := smForBenchmarkThread
	test.BThreadsAtomicWaitFind(b, sm)
}

func BenchmarkThreadsFind(b *testing.B) {
	sm := smForBenchmarkThread
	test.BThreadsFind(b, sm)
}

func BenchmarkThreadsEachShort(b *testing.B) {
	sm := smForBenchmarkThread
	test.BThreadsEachShort(b, sm)
}

func BenchmarkThreadsDelete(b *testing.B) {
	sm := smForBenchmarkThread
	test.BThreadsDelete(b, sm)
}
