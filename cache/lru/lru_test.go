package lru

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

func TestNewCacheReturnsErrorOnInvalidCapacity(t *testing.T) {
	cache, err := NewCache[string, int](0)
	if !errors.Is(err, ErrInvalidCapacity) {
		t.Fatalf("expected ErrInvalidCapacity, got %v", err)
	}
	if cache != nil {
		t.Fatalf("expected nil cache on invalid capacity, got %#v", cache)
	}
}

func TestCacheEvictsLeastRecentlyUsed(t *testing.T) {
	cache, err := NewCache[string, int](2)
	if err != nil {
		t.Fatalf("NewCache returned error: %v", err)
	}
	cache.Put("a", 1)
	cache.Put("b", 2)

	value, ok := cache.Get("a")
	if !ok || value != 1 {
		t.Fatalf("expected key a to exist, got value=%d ok=%v", value, ok)
	}

	cache.Put("c", 3)

	if _, ok := cache.Get("b"); ok {
		t.Fatal("expected key b to be evicted")
	}

	if value, ok := cache.Get("a"); !ok || value != 1 {
		t.Fatalf("expected key a to remain, got value=%d ok=%v", value, ok)
	}
	if value, ok := cache.Get("c"); !ok || value != 3 {
		t.Fatalf("expected key c to exist, got value=%d ok=%v", value, ok)
	}
}

func TestCacheUpdateDoesNotGrowLength(t *testing.T) {
	cache, err := NewCache[string, int](2)
	if err != nil {
		t.Fatalf("NewCache returned error: %v", err)
	}
	cache.Put("a", 1)
	cache.Put("a", 2)

	if cache.Len() != 1 {
		t.Fatalf("expected len 1, got %d", cache.Len())
	}

	value, ok := cache.Get("a")
	if !ok || value != 2 {
		t.Fatalf("expected updated value 2, got value=%d ok=%v", value, ok)
	}
}

func TestCacheDelete(t *testing.T) {
	cache, err := NewCache[string, int](2)
	if err != nil {
		t.Fatalf("NewCache returned error: %v", err)
	}
	cache.Put("a", 1)

	if !cache.Delete("a") {
		t.Fatal("expected delete to succeed")
	}
	if cache.Delete("a") {
		t.Fatal("expected delete on missing key to fail")
	}
	if cache.Len() != 0 {
		t.Fatalf("expected empty cache, got len=%d", cache.Len())
	}
}

func TestCacheConcurrentAccess(t *testing.T) {
	cache, err := NewCache[string, int](8)
	if err != nil {
		t.Fatalf("NewCache returned error: %v", err)
	}
	var wg sync.WaitGroup

	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("key-%d", j%12)
				cache.Put(key, worker+j)
				cache.Get(key)
				if j%5 == 0 {
					cache.Delete(fmt.Sprintf("key-%d", (j+1)%12))
				}
			}
		}(i)
	}

	wg.Wait()

	if cache.Len() > 8 {
		t.Fatalf("expected len <= 8, got %d", cache.Len())
	}
}
