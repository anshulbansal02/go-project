package utils

import "sync"

type KeyMutex struct {
	mutexes sync.Map
}

func (km *KeyMutex) Lock(key string) func() {
	val, _ := km.mutexes.LoadOrStore(key, &sync.Mutex{})
	mtx := val.(*sync.Mutex)
	mtx.Lock()

	return func() { mtx.Unlock() }
}
