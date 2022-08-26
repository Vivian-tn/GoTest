package diplomat

import (
	"sync"
	"time"
)

type Cache struct {
	m sync.Map
}

type memItem struct {
	Value  []*ServiceEntry
	Expire time.Time
}

func NewCache() *Cache {
	return &Cache{}
}

func (m *Cache) Get(key string) ([]*ServiceEntry, bool) {
	v, ok := m.m.Load(key)
	if !ok {
		return nil, true
	}
	item := v.(*memItem)

	return item.Value, item.Expire.Before(time.Now())
}

func (m *Cache) Expire(key string, timeout time.Duration) {
	if value, _ := m.Get(key); value != nil {
		m.m.Store(key, &memItem{Value: value, Expire: time.Now().Add(timeout)})
	}
}

func (m *Cache) Store(key string, value []*ServiceEntry, timeout time.Duration) {
	m.m.Store(key, &memItem{Value: value, Expire: time.Now().Add(timeout)})
}

func (m *Cache) Discard(key string, discarded *ServiceEntry) {
	v, ok := m.m.Load(key)
	if !ok {
		return
	}
	item := v.(*memItem)

	newEntries := make([]*ServiceEntry, 0)
	for _, entry := range item.Value {
		if !(entry.Host == discarded.Host && entry.Port == discarded.Port) {
			newEntries = append(newEntries, entry)
		}
	}
	m.m.Store(key, &memItem{Value: newEntries, Expire: item.Expire})
}

func (m *Cache) Delete(key string) {
	m.m.Delete(key)
}

func (m *Cache) Has(key string) bool {
	_, ok := m.m.Load(key)
	return ok
}
