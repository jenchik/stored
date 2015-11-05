package multimap

import (
	"testing"

	"github.com/jenchik/stored/api"
	"github.com/jenchik/stored/test"
)

var smForBenchmark api.StoredMap
var smForBenchmarkThread api.StoredMap

func init() {
	smForBenchmark = newTest()
	err := test.InserterBasic(smForBenchmark, "Benchmark")
	if err != nil {
		panic(err.Error())
	}
	err = testWaitN(smForBenchmark, test.CntWorks*test.CntItems)
	if err != nil {
		panic(err.Error())
	}
	smForBenchmarkThread = newTest()
	err = test.InserterBasic(smForBenchmarkThread, "BenchmarkThreads")
	if err != nil {
		panic(err.Error())
	}
	err = testWaitN(smForBenchmarkThread, test.CntWorks*test.CntItems)
	if err != nil {
		panic(err.Error())
	}
}

func BenchmarkInsert(b *testing.B) {
	sm := newTest()
	test.BInsert(b, sm)
}

func BenchmarkAtomicUpdate(b *testing.B) {
	sm := newTest()
	test.BAtomicUpdate(b, sm)
}

func BenchmarkAtomicWaitUpdate(b *testing.B) {
	sm := newTest()
	test.BAtomicWaitUpdate(b, sm)
}

func BenchmarkUpdate(b *testing.B) {
	sm := newTest()
	test.BUpdate(b, sm)
}

func BenchmarkAtomicComplex(b *testing.B) {
	sm := newTest()
	test.BAtomicComplex(b, sm)
}

func BenchmarkAtomicWaitComplex(b *testing.B) {
	sm := newTest()
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

func BenchmarkAtomicEachN(b *testing.B) {
	sm := smForBenchmark
	test.BAtomicEachN(b, sm)
}

// TODO
/*
func BenchmarkAtomicEachShort(b *testing.B) {
	sm := smForBenchmark
	test.BAtomicEachShort(b, sm)
}
*/

func BenchmarkEachN(b *testing.B) {
	sm := smForBenchmark
	test.BEachN(b, sm)
}

// TODO
/*
func BenchmarkEachShort(b *testing.B) {
	sm := smForBenchmark
	test.BEachShort(b, sm)
}
*/

func BenchmarkDelete(b *testing.B) {
	sm := smForBenchmark
	test.BDelete(b, sm)
}

func BenchmarkThreadsInsert(b *testing.B) {
	sm := newTest()
	test.BThreadsInsert(b, sm)
}

// TODO
/*
func BenchmarkThreadsAtomicUpdate(b *testing.B) {
	sm := newTest()
	test.BThreadsAtomicUpdate(b, sm)
}

func BenchmarkThreadsAtomicWaitUpdate(b *testing.B) {
	sm := newTest()
	test.BThreadsAtomicWaitUpdate(b, sm)
}

func BenchmarkThreadsUpdate(b *testing.B) {
	sm := newTest()
	test.BThreadsUpdate(b, sm)
}

func BenchmarkThreadsAtomicComplex(b *testing.B) {
	sm := newTest()
	test.BThreadsAtomicComplex(b, sm)
}

func BenchmarkThreadsAtomicWaitComplex(b *testing.B) {
	sm := newTest()
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
*/

func BenchmarkThreadsDelete(b *testing.B) {
	sm := smForBenchmarkThread
	test.BThreadsDelete(b, sm)
}
