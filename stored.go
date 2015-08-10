package stored

import (
	"github.com/jenchik/spctools"
	"github.com/jenchik/stored/api"
	"math"
	"time"
)

var (
	TTL        int = 100 // in seconds
	MaxItems   int = 100
	IntervalGC int = 100 // in seconds
)

type Cacher interface {
	GetItem(key string, fback api.GetterFunc) (interface{}, error)
	ResetItem(key string)
	Stats() map[string]uint64
	//SetOptions(Options) // TODO
}

type GCer interface {
	timeoutGC() time.Duration
	GC()
}

type GCCacher interface {
	GCer
	Cacher
}

type Options struct {
	MaxItems   int
	TTL        int // in seconds
	TimeoutGC  int // in seconds
	DisabledGC bool
}

var saneMaxItems, saneTimeoutGC, saneTTL func(int) int

func init() {
	saneMaxItems = spctools.MakeBoundedIntFunc(10, math.MaxInt32, &MaxItems)
	saneTimeoutGC = spctools.MakeBoundedIntFunc(4, 86400, &IntervalGC) // max = 24 hours
	saneTTL = spctools.MakeBoundedIntFunc(2, 86400, &TTL)              // max = 24 hours
}

func gc(c GCCacher) {
	for {
		<-time.After(c.timeoutGC() * time.Second)
		c.GC()
	}
}
