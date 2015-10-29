package multimap

import (
	"testing"

	"github.com/jenchik/stored/api"
	smap "github.com/jenchik/stored/safemap"
	"github.com/jenchik/stored/test"
)

// WARNING! Setted 'unsafe' mode
const testUnsafeMultiN = -4

var smUnsafeForBenchmark api.StoredMap

func init() {
	smUnsafeForBenchmark = newUnsafeTest()
	err := test.InserterBasic(smUnsafeForBenchmark, "BenchmarkUnsafe")
	if err != nil {
		panic(err.Error())
	}
}

func newUnsafeTest() api.StoredMap {
	return New(testUnsafeMultiN, smap.New().(api.StoredCopier))
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
	sm := newUnsafeTest()
	test.BInsert(b, sm)
}

func BenchmarkUnsafeAtomicUpdate(b *testing.B) {
	sm := newUnsafeTest()
	test.BAtomicUpdate(b, sm)
}

func BenchmarkUnsafeAtomicWaitUpdate(b *testing.B) {
	sm := newUnsafeTest()
	test.BAtomicWaitUpdate(b, sm)
}

func BenchmarkUnsafeUpdate(b *testing.B) {
	sm := newUnsafeTest()
	test.BUpdate(b, sm)
}

func BenchmarkUnsafeAtomicComplex(b *testing.B) {
	sm := newUnsafeTest()
	test.BAtomicComplex(b, sm)
}

func BenchmarkUnsafeAtomicWaitComplex(b *testing.B) {
	sm := newUnsafeTest()
	test.BAtomicWaitComplex(b, sm)
}

func BenchmarkUnsafeAtomicFind(b *testing.B) {
	sm := smUnsafeForBenchmark
	test.BAtomicFind(b, sm)
}

func BenchmarkUnsafeAtomicWaitFind(b *testing.B) {
	sm := smUnsafeForBenchmark
	test.BAtomicWaitFind(b, sm)
}

func BenchmarkUnsafeFind(b *testing.B) {
	sm := smUnsafeForBenchmark
	test.BFind(b, sm)
}

// TODO
/*
func BenchmarkUnsafeEachFullCicle(b *testing.B) {
	sm := smUnsafeForBenchmark
    test.BEachFullCicle(b, sm)
}

func BenchmarkUnsafeEachShort(b *testing.B) {
	sm := smUnsafeForBenchmark
    test.BEachShort(b, sm)
}
*/

func BenchmarkUnsafeDelete(b *testing.B) {
	sm := smUnsafeForBenchmark
	test.BDelete(b, sm)
}
