package storage

import (
	"sync"
	"time"
)

type Item struct {
	Value string
	ExpiresAt int64
}

type DB struct {
	mu sync.RWMutex
	data map[string]Item
}

func NewDB() *DB {
	return &DB{
		data: make(map[string]Item),
	}
}

func (db *DB) Set(key string, value string, timestamp int64) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key] = Item {
		Value: value,
		ExpiresAt: timestamp,
	}
}

func (db *DB) Get(key string) (string, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	item, ok := db.data[key]
	return item.Value, ok
}

func (db *DB) Del(key string) int {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	if _, ok := db.data[key]; ok {
		delete(db.data, key)
		return 1
	}
	return 0
}

func (db *DB) StartJanitor(interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for {
			<-ticker.C

			db.DeleteExpired()
		}
	}()
}

func (db *DB) DeleteExpired() {
	db.mu.Lock()
	defer db.mu.Unlock()

	now := time.Now().UnixNano()
	for key, item := range db.data {
		if item.ExpiresAt > 0 && now > item.ExpiresAt {
			delete(db.data, key)
		}
	}
}