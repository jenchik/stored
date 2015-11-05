package api

type AtomicFunc func(Mapper)
type ForeachFunc func(Mapper) // deprecated
type UpdateFunc func(interface{}, bool) interface{}
type GetterFunc func(string) (interface{}, error)

type Iterator interface {
	Done() <-chan struct{}
	Next() bool
	Stop()
}

type Mapper interface {
	Find(key string) (value interface{}, found bool)
	Key() string
	SetKey(key string)
	Value() interface{}
	Delete()
	Update(value interface{})
	Len() int
	Lock()
	Unlock()
	Stop()
	Next() bool
	Clear()
	Close() // do not use!
}

type StoredCopier interface {
	Copy() StoredMap
}

type StoredMap interface {
	Delete(string)
	Find(string) (interface{}, bool)
	Insert(string, interface{})
	Atomic(AtomicFunc)
	AtomicWait(AtomicFunc)
	Len() int
	Each(ForeachFunc) // deprecated
	Update(string, UpdateFunc)
}
