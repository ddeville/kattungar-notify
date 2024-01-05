package gcal

import (
	"container/list"
)

type LRUCache struct {
	capacity int
	queue    *list.List
	items    map[string]*list.Element
}

func NewLRUCache(capacity int) LRUCache {
	return LRUCache{
		capacity: capacity,
		queue:    list.New(),
		items:    make(map[string]*list.Element),
	}
}

func (c *LRUCache) Add(value string) {
	if elem, found := c.items[value]; found {
		c.queue.MoveToFront(elem)
		return
	}

	if c.queue.Len() == c.capacity {
		oldest := c.queue.Back()
		if oldest != nil {
			delete(c.items, oldest.Value.(string))
			c.queue.Remove(oldest)
		}
	}

	elem := c.queue.PushFront(value)
	c.items[value] = elem
}

func (c *LRUCache) Contains(value string) bool {
	_, found := c.items[value]
	return found
}
