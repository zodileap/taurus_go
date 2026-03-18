package lru

import (
	"container/list"
	"errors"
	"fmt"
	"sync"
)

var ErrInvalidCapacity = errors.New("lru cache capacity must be positive")

type entry[K comparable, V any] struct {
	key   K
	value V
}

// Cache 提供一个并发安全的 LRU 缓存。
type Cache[K comparable, V any] struct {
	capacity int
	items    map[K]*list.Element
	order    *list.List
	mu       sync.RWMutex
}

// NewCache 创建一个指定容量的缓存。
func NewCache[K comparable, V any](capacity int) (*Cache[K, V], error) {
	if capacity <= 0 {
		return nil, fmt.Errorf("%w: %d", ErrInvalidCapacity, capacity)
	}

	return &Cache[K, V]{
		capacity: capacity,
		items:    make(map[K]*list.Element, capacity),
		order:    list.New(),
	}, nil
}

// Get 返回命中的值，并将该键提升为最近使用。
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var zero V
	element, ok := c.items[key]
	if !ok {
		return zero, false
	}

	c.order.MoveToFront(element)
	return element.Value.(*entry[K, V]).value, true
}

// Put 写入键值，并在容量超限时淘汰最近最少使用项。
func (c *Cache[K, V]) Put(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, ok := c.items[key]; ok {
		c.order.MoveToFront(element)
		element.Value.(*entry[K, V]).value = value
		return
	}

	element := c.order.PushFront(&entry[K, V]{
		key:   key,
		value: value,
	})
	c.items[key] = element

	if c.order.Len() <= c.capacity {
		return
	}

	oldest := c.order.Back()
	if oldest == nil {
		return
	}

	c.order.Remove(oldest)
	delete(c.items, oldest.Value.(*entry[K, V]).key)
}

// Delete 删除指定键。
func (c *Cache[K, V]) Delete(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	element, ok := c.items[key]
	if !ok {
		return false
	}

	c.order.Remove(element)
	delete(c.items, key)
	return true
}

// Len 返回当前缓存条目数。
func (c *Cache[K, V]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}
