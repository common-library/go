// Package lock provides lock implementations.
package lock

import "sync"

// Mutex is a struct that provides mutex related methods.
type Mutex struct {
	mutex sync.Mutex
}

// Lock is lock.
//
// ex) mutex.Lock()
func (m *Mutex) Lock() {
	m.mutex.Lock()
}

// TryLock attempts a lock and returns whether it was successful or not.
//
// ex) result := mutex.TryLock()
func (m *Mutex) TryLock() bool {
	return m.mutex.TryLock()
}

// Unlock is unlock.
//
// ex) mutex.Unlock()
func (m *Mutex) Unlock() {
	m.mutex.Unlock()
}

// MutexByKey is a struct that provides mutex-related methods for each key.
type MutexByKey struct {
	mutexs sync.Map
}

// Lock locks the mutex corresponding to the key.
//
// ex) mutexs.Lock(key)
func (mbk *MutexByKey) Lock(key any) {
	mutex, _ := mbk.mutexs.LoadOrStore(key, &Mutex{})

	mutex.(*Mutex).Lock()
}

// TryLock attempts to lock the mutex corresponding to the key and returns whether it was successful or not.
//
// ex) result := mutexs.TryLock(key)
func (mbk *MutexByKey) TryLock(key any) bool {
	mutex, _ := mbk.mutexs.LoadOrStore(key, &Mutex{})

	return mutex.(*Mutex).TryLock()
}

// Unlock unlocks the mutex corresponding to the key.
//
// ex) mutexs.Unlock(key)
func (mbk *MutexByKey) Unlock(key any) {
	mutex, _ := mbk.mutexs.LoadOrStore(key, &Mutex{})

	mutex.(*Mutex).Unlock()
}

// UnlockAndDelete unlocks and deletes the mutex corresponding to the key.
//
// ex) mutexs.UnlockAndDelete(key)
func (mbk *MutexByKey) UnlockAndDelete(key any) {
	if mutex, loaded := mbk.mutexs.LoadAndDelete(key); loaded {
		mutex.(*Mutex).Unlock()
	} else {
		(&Mutex{}).Unlock()
	}
}

// Delete deletes the mutex corresponding to the key.
//
// ex) mutexs.Delete(key)
func (mbk *MutexByKey) Delete(key any) {
	mbk.mutexs.Delete(key)
}
