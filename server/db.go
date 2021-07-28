package server

import (
	"sync"
)

type DB struct {
	mutex sync.RWMutex
	store map[string]string
}

func NewDB() *DB {
	return &DB{
		store: make(map[string]string),
	}
}

func (db *DB) Get(key string) (string, bool) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	value, err := db.store[key]
	if !err {
		return "", false
	}
	return value, true
}

func (db *DB) Put(key string, value string) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.store[key] = value
}
