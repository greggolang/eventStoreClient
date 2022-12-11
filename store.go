package eventeur

import (
	"sync"
	"sync/atomic"
)

type Store[T any] struct {
	// mu is protecting the internal store being modified in parallel and
	// inserting different items at the same idx.
	//
	// It also controls the read operations on internal store. For example,
	// When a new item is being inserted, all parallel reads will wait until
	// the `Put` operation function is completed.
	mu sync.RWMutex

	// Internal store for items of type T (this can be Client, Topic or any other
	// type).
	store map[uint64]*T

	// idx is atomically incremented for every newly inserted item
	// into the store.
	//
	// Atomic operation allows idx to be incremented in parallel without
	// a fear of it being modified at the same time.
	idx uint64
}

func NewStore[T any]() *Store[T] {
	return &Store[T]{
		store: make(map[uint64]*T),
	}
}

func (s *Store[T]) Put(item *T) uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.newID()
	s.store[id] = item
	return id
}

// Update returns true if the item at given id exists and the update succeeds.
// If an item does not exist, it returns false without modifying internal store.
func (s *Store[T]) Update(id uint64, item *T) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if an item exists before updating it.
	if _, ok := s.store[id]; !ok {
		return false
	}

	s.store[id] = item
	return true
}

// Get return an item and `true` as a second return parameter if the item with a
// provided id exists. Otherwise the item will be `nil` and second param `false`.
func (s *Store[T]) Get(id uint64) (*T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.store[id]
	return item, ok
}

func (s *Store[T]) Delete(id uint64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if an item exists before deleting it.
	if _, ok := s.store[id]; !ok {
		return false
	}

	delete(s.store, id)
	return true
}

func (s *Store[T]) ForEach(f func(uint64, *T)) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for id, v := range s.store {
		f(id, v)
	}
}

func (s *Store[T]) newID() uint64 {
	atomic.AddUint64(&s.idx, 1)
	return s.idx
}
