package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type ListValueWrapper struct {
	key   Key
	value interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	wv := wrapValue(key, value)

	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		item.Value = wv
		return true
	}

	if c.queue.Len() == c.capacity {
		item := c.queue.Back()
		c.queue.Remove(item)
		delete(c.items, item.Value.(ListValueWrapper).key)
	}

	item := c.queue.PushFront(wv)
	c.items[key] = item

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		return item.Value.(ListValueWrapper).value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func wrapValue(key Key, value interface{}) ListValueWrapper {
	return ListValueWrapper{
		key:   key,
		value: value,
	}
}
