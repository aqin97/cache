package cache

import (
	"log"
	"sync"
)

type Cache interface {
	Set(string, []byte) error
	Get(string) ([]byte, error)
	Delete(string) error
	GetStat() Stat
}

type Stat struct {
	Count     int64
	KeySize   int64
	ValueSize int64
}

func (s *Stat) add(k string, v []byte) {
	s.Count += 1
	s.KeySize += int64(len(k))
	s.ValueSize += int64(len(v))
}

func (s *Stat) del(k string, v []byte) {
	s.Count -= 1
	s.KeySize -= int64(len(k))
	s.ValueSize -= int64(len(v))
}

func New(typ string) Cache {
	var c Cache
	if typ == "inmemory" {
		c = newInMemoryCache()
	}
	if c == nil {
		panic("unknow cache type " + typ)
	}
	log.Println(typ, "ready to serve")
	return c
}

type inMemoryCache struct {
	c     map[string][]byte
	mutex sync.RWMutex
	Stat
}

func (i *inMemoryCache) Set(k string, v []byte) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	tmp, exist := i.c[k]
	if exist {
		i.del(k, tmp)
	}
	i.c[k] = v
	i.add(k, v)
	return nil
}

func (i *inMemoryCache) Get(k string) ([]byte, error) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	return i.c[k], nil
}

func (i *inMemoryCache) Delete(k string) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	v, exist := i.c[k]
	if exist {
		delete(i.c, k)
		i.del(k, v)
	}
	return nil
}

func (i *inMemoryCache) GetStat() Stat {
	return i.Stat
}

func newInMemoryCache() Cache {
	return &inMemoryCache{make(map[string][]byte), sync.RWMutex{}, Stat{}}
}
