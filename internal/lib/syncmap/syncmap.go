package syncmap

import "sync"

type SyncMap[K comparable, V any] struct {
	m  map[K]V
	mx sync.Mutex
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		m:  make(map[K]V),
		mx: sync.Mutex{},
	}
}

func (s *SyncMap[K, V]) Store(key K, value V) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.m[key] = value
}

func (s *SyncMap[K, V]) Load(key K) (V, bool) {
	s.mx.Lock()
	defer s.mx.Unlock()
	value, ok := s.m[key]
	return value, ok
}

func (s *SyncMap[K, V]) Delete(key K) {
	s.mx.Lock()
	defer s.mx.Unlock()
	delete(s.m, key)
}
