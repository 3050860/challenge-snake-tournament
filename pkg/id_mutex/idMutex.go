package idMutex

import (
	"fmt"
	"sync"
)

const HeartbeatURL = "/userapi/heartbeat"

type IdMutex struct {
	mutex   sync.Mutex
	mutexes map[string]*mutexEntry
}

type mutexEntry struct {
	idMutex *IdMutex
	mutex   sync.Mutex
	count   int
	key     string
}

func New() *IdMutex {
	return &IdMutex{mutexes: make(map[string]*mutexEntry)}
}

type Unlocker interface {
	Unlock()
}

func (m *IdMutex) Lock(key string) Unlocker {
	m.mutex.Lock()
	mutex, ok := m.mutexes[key]
	if !ok {
		mutex = &mutexEntry{idMutex: m, key: key}
		m.mutexes[key] = mutex
	}
	mutex.count++
	m.mutex.Unlock()

	mutex.mutex.Lock()

	return mutex
}

func (me *mutexEntry) Unlock() {
	idMutex := me.idMutex

	idMutex.mutex.Lock()
	mutex, ok := idMutex.mutexes[me.key]
	if !ok { // entry must exist
		idMutex.mutex.Unlock()
		panic(fmt.Errorf("Unlock requested for key=%v but no entry found", me.key))
	}
	mutex.count--        // ref count
	if mutex.count < 1 { // if it hits zero then we own it and remove from map
		delete(idMutex.mutexes, me.key)
	}

	idMutex.mutex.Unlock()

	mutex.mutex.Unlock()
}
