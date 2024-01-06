package gcal

import (
	"testing"
)

func TestLRUCache_AddAndContains(t *testing.T) {
	cache := NewLRUCache(2)
	cache.Add("item1")
	cache.Add("item2")

	if !cache.Contains("item1") {
		t.Errorf("Cache should contain 'item1'")
	}

	if !cache.Contains("item2") {
		t.Errorf("Cache should contain 'item2'")
	}
}

func TestLRUCache_Capacity(t *testing.T) {
	cache := NewLRUCache(2)
	cache.Add("item1")
	cache.Add("item2")
	cache.Add("item3") // This should evict "item1"

	if cache.Contains("item1") {
		t.Error("Cache should not contain 'item1'")
	}

	if !cache.Contains("item3") {
		t.Errorf("Cache should contain 'item3'")
	}
}

func TestLRUCache_UpdateOrderOnAdd(t *testing.T) {
	cache := NewLRUCache(2)
	cache.Add("item1")
	cache.Add("item2")
	cache.Add("item1") // This should make item2 the least recently used
	cache.Add("item3") // This should evict "item2"

	if cache.Contains("item2") {
		t.Error("Cache should not contain 'item2'")
	}

	if !cache.Contains("item1") {
		t.Errorf("Cache should contain 'item1'")
	}
}
