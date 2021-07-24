package db

import (
	"sync"
)

type DB struct {
	mutex sync.Mutex
	store map[string]string
}

func NewDB() *DB {
	return &DB{
		store: make(map[string]string),
	}
}

func (db *DB) Get(key string) (string, bool) {
	db.mutex.Lock()
	value, err := db.store[key]
	if !err {
		return "", false
	}
	return value, true
}

func (db *DB) Put(key string, value string) {
	db.mutex.Lock()
	db.store[key] = value
}
