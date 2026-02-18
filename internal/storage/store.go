package storage

import (
	"sync"
)

type DB struct {
	mu sync.RWMutex
	data map[string]string
}

func NewDB() *DB {
	return &DB{
		data: make(map[string]string),
	}
}

func (db *DB) Set(key string, value string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key] = value
}

func (db *DB) Get(key string) (string, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	val, ok := db.data[key]
	return val, ok
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