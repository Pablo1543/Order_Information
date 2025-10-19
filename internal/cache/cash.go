package cache

import (
	"Order_Information/internal/db"
	"context"
	"encoding/json"
	"sync"
)

type Cache struct {
	mu sync.RWMutex
	m  map[string]json.RawMessage
}

func NewCache() *Cache {
	return &Cache{m: make(map[string]json.RawMessage)}
}

func (c *Cache) Get(id string) (json.RawMessage, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.m[id]
	return v, ok
}

func (c *Cache) Set(id string, raw json.RawMessage) {
	c.mu.Lock()
	c.m[id] = raw
	c.mu.Unlock()
}

func (c *Cache) LoadFromDB(database *db.DB) error {
	ctx := context.Background()
	raws, err := database.LoadAllOrders(ctx)
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for id, raw := range raws {
		c.m[id] = json.RawMessage(raw)
	}
	return nil
}

func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.m)
}

func (c *Cache) GetAllKeys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.m))
	for k := range c.m {
		keys = append(keys, k)
	}
	return keys
}
