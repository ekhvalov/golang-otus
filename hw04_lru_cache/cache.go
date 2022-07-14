package hw04lrucache

import "sync"

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
	mutex    sync.Mutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	lItem, isExists := l.items[key]
	cItem := cacheItem{key: key, value: value}
	if isExists {
		lItem.Value = cItem
		l.queue.MoveToFront(lItem)
		return true
	}
	if l.queue.Len() == l.capacity {
		l.removeLastUsedItem()
	}
	lItem = l.queue.PushFront(cItem)
	l.items[key] = lItem
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	item, isExists := l.items[key]
	if isExists {
		l.queue.MoveToFront(item)
		return item.Value.(cacheItem).value, isExists
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

func (l *lruCache) removeLastUsedItem() {
	lastItem := l.queue.Back()
	lastCacheItem := lastItem.Value.(cacheItem)
	l.queue.Remove(lastItem)
	delete(l.items, lastCacheItem.key)
}
