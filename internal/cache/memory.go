package cache

import (
	"context"
	"sync"
	"time"
)

type memoryItem struct {
	value      []byte
	expiration time.Time
}

type memory struct {
	storage   map[string]*memoryItem
	storageMu sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewMemoryCache provides cache in memory
func NewMemoryCache(ctx context.Context) Cache {
	c, cancel := context.WithCancel(ctx)
	m := &memory{
		storage: make(map[string]*memoryItem),
		ctx:     c,
		cancel:  cancel,
	}
	m.cleanUp(ctx)
	return m
}

func (m *memory) cleanUp(ctx context.Context) {
	ticker := time.NewTimer(time.Second)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-m.ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			keys := make([]string, 0)
			t := time.Now().UTC()
			m.storageMu.RLock()

			for key, value := range m.storage {
				if t.After(value.expiration) {
					keys = append(keys, key)
				}
			}
			m.storageMu.RUnlock()
			if len(keys) > 0 {
				m.storageMu.Lock()
				for _, key := range keys {
					delete(m.storage, key)
				}
				m.storageMu.Unlock()
			}
		}
	}
}

// Get data from memory
func (m *memory) Get(key string) ([]byte, error) {
	m.storageMu.RLock()
	defer m.storageMu.RUnlock()
	data, ok := m.storage[key]
	if !ok {
		return nil, nil
	}
	if time.Now().UTC().After(data.expiration) {
		return nil, nil
	}
	return data.value, nil
}

// Set data to memory
func (m *memory) Set(key string, value []byte, expiration time.Duration) error {
	m.storageMu.Lock()
	defer m.storageMu.Unlock()
	m.storage[key] = &memoryItem{value: value, expiration: time.Now().UTC().Add(expiration)}
	return nil
}

// Delete data from memory
func (m *memory) Delete(key string) error {
	m.storageMu.Lock()
	defer m.storageMu.Unlock()
	delete(m.storage, key)
	return nil
}

// Close cache
func (m *memory) Close() error {
	m.cancel()
	return nil
}
