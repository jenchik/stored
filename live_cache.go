package stored

import (
	"github.com/jenchik/stored/api"
)

type LiveCache struct {
	*Cache
}

func NewLiveCache(opt Options, sm api.StoredMap) GCCacher {
	lc := &LiveCache{NewCache(opt, sm)}
	return lc
}
