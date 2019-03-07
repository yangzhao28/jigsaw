package vars

import (
	"sync"
)

type Storage interface {
	Persist() error
	Restore() error
	Put(string, interface{})
	Get(string) interface{}
	Del(string) interface{}
	Foreach(func(string, interface{}) error) error
}

type MemoryStorage struct {
	kv map[string]interface{}
	l  sync.Mutex
}

func NewMemoryStorage() Storage {
	return &MemoryStorage{
		kv: make(map[string]interface{}),
	}
}

func (s *MemoryStorage) Persist() error {
	return nil
}

func (s *MemoryStorage) Restore() error {
	return nil
}

func (s *MemoryStorage) Get(key string) interface{} {
	s.l.Lock()
	defer s.l.Unlock()
	return s.kv[key]
}

func (s *MemoryStorage) Put(key string, val interface{}) {
	s.l.Lock()
	defer s.l.Unlock()
	s.kv[key] = val
}

func (s *MemoryStorage) Del(key string) interface{} {
	s.l.Lock()
	defer s.l.Unlock()
	v := s.kv[key]
	delete(s.kv, key)
	return v
}

func (s *MemoryStorage) Foreach(f func(string, interface{}) error) error {
	for k, v := range s.kv {
		err := f(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
