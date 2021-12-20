package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

// Добавить значение в кэш по ключу.
func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	cItem := cacheItem{key, value}
	item, ok := c.items[key]
	// при наличии элемента, обновляем значение и перемещаем в начало
	if ok {
		item.Value = cItem
		c.queue.MoveToFront(item)
		return true
	}

	// проверка длины очереди
	if c.queue.Len() >= c.capacity {
		last := c.queue.Back()

		if v, ok := last.Value.(cacheItem); ok {
			delete(c.items, v.key)
		}
		c.queue.Remove(last)
	}

	// добавляем элемент в очередь и указатель в словарь
	c.items[key] = c.queue.PushFront(cItem)

	return false
}

// Получить значение из кэша по ключу.
func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(item)
		switch v := item.Value.(type) {
		case cacheItem:
			return v.value, ok
		default:
			return item, ok
		}
	}
	return item, ok
}

// Очистить кэш.
func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
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
