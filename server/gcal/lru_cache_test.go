package gcal

import (
	"testing"
)

func TestLRUCache_SetAndGet(t *testing.T) {
	cache := NewLRUCache(2)
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	if val, ok := cache.Get("key1"); !ok || val != "value1" {
		t.Errorf("Got %s; want %s", val, "value1")
	}

	if val, ok := cache.Get("key2"); !ok || val != "value2" {
		t.Errorf("Got %s; want %s", val, "value2")
	}
}

func TestLRUCache_Capacity(t *testing.T) {
	cache := NewLRUCache(2)
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3") // This should evict "key1"

	if _, ok := cache.Get("key1"); ok {
		t.Error("Expected key1 to be evicted")
	}

	if val, ok := cache.Get("key3"); !ok || val != "value3" {
		t.Errorf("Got %s; want %s", val, "value3")
	}
}

func TestLRUCache_UpdateOrderOnGet(t *testing.T) {
	cache := NewLRUCache(2)
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Get("key1")           // This should make key2 the least recently used
	cache.Set("key3", "value3") // This should evict "key2"

	if _, ok := cache.Get("key2"); ok {
		t.Error("Expected key2 to be evicted")
	}

	if val, ok := cache.Get("key1"); !ok || val != "value1" {
		t.Errorf("Got %s; want %s", val, "value1")
	}
}
