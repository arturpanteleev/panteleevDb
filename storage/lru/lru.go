package lru

import (
	"container/list"
	"errors"
)

type (
	LRU struct {
		size      int
		evictList *list.List
		items     map[interface{}]*list.Element
	}
	// entry used to store value in evictList
	entry struct {
		key   interface{}
		value interface{}
	}
)

// New initialized a new LRU with fixed size
func New(size int) (*LRU, error) {
	if size <= 0 {
		return nil, errors.New("Size must be greater than 0")
	}
	c := &LRU{
		size:      size,
		evictList: list.New(),
		items:     make(map[interface{}]*list.Element),
	}
	return c, nil
}

// Len returns the number of items in cache
func (l *LRU) Len() int {
	return l.evictList.Len()
}

// Add adds a value to the cache. Return true if eviction occured
func (l *LRU) Add(key, value interface{}) bool {
	if ent, ok := l.items[key]; ok {
		l.evictList.MoveToFront(ent)
		ent.Value.(*entry).value = value
		return false
	}
	ent := &entry{key, value}
	entry := l.evictList.PushFront(ent)
	l.items[key] = entry
	evict := l.evictList.Len() > l.size
	if evict {
		l.removeOldest()
	}
	return evict
}

func (l *LRU) removeOldest() {
	ent := l.evictList.Back()
	if ent != nil {
		l.removeElement(ent)
	}
}
func (l *LRU) removeElement(e *list.Element) {
	l.evictList.Remove(e)
	kv := e.Value.(*entry)
	delete(l.items, kv.key)
}

// Purge completely clears cache
func (l *LRU) Purge() {
	for k := range l.items {
		delete(l.items, k)
	}
	l.evictList.Init()
}

// Get looks up a key's value from the cache
func (l *LRU) Get(key interface{}) (value interface{}, ok bool) {
	if ent, ok := l.items[key]; ok {
		l.evictList.MoveToFront(ent)
		return ent.Value.(*entry).value, true
	}
	return
}

// Contains check if key is in cache without updating
// recent-ness or deleting it for being state.
func (l *LRU) Contains(key interface{}) (ok bool) {
	_, ok = l.items[key]
	return ok
}

// RemoveOldest removes oldest item from cache
func (l *LRU) RemoveOldest() (interface{}, interface{}, bool) {
	ent := l.evictList.Back()
	if ent != nil {
		l.removeElement(ent)
		kv := ent.Value.(*entry)
		return kv.key, kv.value, true
	}
	return nil, nil, false
}

// GetOldest returns oldest item from cache
func (l *LRU) GetOldest() (interface{}, interface{}, bool) {
	ent := l.evictList.Back()
	if ent != nil {
		kv := ent.Value.(*entry)
		return kv.key, kv.value, true
	}
	return nil, nil, false
}
