package cache

import (
	"sync"
	"time"
)

type Cache struct {
	sync.RWMutex
	data map[string]interface{}
}

func New() *Cache {
	return &Cache{
		data: make(map[string]interface{}),
	}
}

// Remove is a func to remove specific key from cache
func (c *Cache) Remove(key string) {
	c.Lock()
	delete(c.data, key)
	c.Unlock()
}

// Keys is a func to get all available keys
func (c *Cache) Keys() []string {
	c.RLock()
	defer c.RUnlock()

	result := make([]string, 0, len(c.data))
	for key := range c.data {
		result = append(result, key)
	}

	return result
}

// SetTTL is a func to set TTL for existing key
func (c *Cache) SetTTL(key string, ttl time.Duration) error {
	if ttl < 0 {
		return ErrorInvalidTTL
	}

	if ttl == 0 {
		return nil
	}

	_, ok := c.data[key]
	if !ok {
		return ErrorKeyNotFound
	}

	time.AfterFunc(ttl, func() {
		c.Remove(key)
	})

	return nil
}

// Set is a func to set String value of a key
func (c *Cache) Set(key string, value string) {
	c.Lock()
	c.data[key] = String(value)
	c.Unlock()
}

// Get is a func to get String value of a key
func (c *Cache) Get(key string) (string, error) {
	c.RLock()
	defer c.RUnlock()

	item, ok := c.data[key]
	if !ok {
		return "", ErrorKeyNotFound
	}

	str, ok := item.(String)
	if !ok {
		return "", ErrorWrongType
	}

	return string(str), nil
}

// Push is a func to append values to a list
func (c *Cache) Push(key string, values ...string) error {
	c.Lock()
	defer c.Unlock()

	if item, ok := c.data[key]; ok {
		list, ok := item.(List)
		if !ok {
			return ErrorWrongType
		}
		c.data[key] = append(list, values...)
	} else {
		c.data[key] = List(values)
	}

	return nil
}

// Pop is a func to remove and get first item in a list
func (c *Cache) Pop(key string) (string, error) {
	c.Lock()
	defer c.Unlock()

	item, ok := c.data[key]
	if !ok {
		return "", ErrorKeyNotFound
	}

	list, ok := item.(List)
	if !ok {
		return "", ErrorWrongType
	}

	if len(list) == 0 {
		return "", ErrorNoItems
	}

	x := list[0]
	c.data[key] = list[1:]
	return x, nil
}

// Hset is a func to set the string value of a dict field
func (c *Cache) Hset(key, field, value string) error {
	c.Lock()
	defer c.Unlock()

	if item, ok := c.data[key]; ok {
		dict, ok := item.(Dict)
		if !ok {
			return ErrorWrongType
		}
		dict[field] = value
	} else {
		c.data[key] = Dict(map[string]string{field: value})
	}

	return nil
}

// Hget is a func to get value of a dict field
func (c *Cache) Hget(key, field string) (string, error) {
	c.RLock()
	defer c.RUnlock()

	item, ok := c.data[key]
	if !ok {
		return "", ErrorKeyNotFound
	}

	dict, ok := item.(Dict)
	if !ok {
		return "", ErrorWrongType
	}

	return dict[field], nil
}
