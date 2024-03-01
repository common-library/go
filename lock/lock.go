package lock

import "sync"

// Mutex is a struct that provides mutex related methods.
type Mutex struct {
	mutex sync.Mutex
}

// Lock is lock.
//
// ex) mutex.Lock()
func (this *Mutex) Lock() {
	this.mutex.Lock()
}

// TryLock attempts a lock and returns whether it was successful or not.
//
// ex) result := mutex.TryLock()
func (this *Mutex) TryLock() bool {
	return this.mutex.TryLock()
}

// Unlock is unlock.
//
// ex) mutex.Unlock()
func (this *Mutex) Unlock() {
	this.mutex.Unlock()
}

// MutexByKey is a struct that provides mutex-related methods for each key.
type MutexByKey struct {
	mutexs sync.Map
}

// Lock locks the mutex corresponding to the key.
//
// ex) mutexs.Lock(key)
func (this *MutexByKey) Lock(key any) {
	mutex, _ := this.mutexs.LoadOrStore(key, &Mutex{})

	mutex.(*Mutex).Lock()
}

// TryLock attempts to lock the mutex corresponding to the key and returns whether it was successful or not.
//
// ex) result := mutexs.TryLock(key)
func (this *MutexByKey) TryLock(key any) bool {
	mutex, _ := this.mutexs.LoadOrStore(key, &Mutex{})

	return mutex.(*Mutex).TryLock()
}

// Unlock unlocks the mutex corresponding to the key.
//
// ex) mutexs.Unlock(key)
func (this *MutexByKey) Unlock(key any) {
	mutex, _ := this.mutexs.LoadOrStore(key, &Mutex{})

	mutex.(*Mutex).Unlock()
}

// UnlockAndDelete unlocks and deletes the mutex corresponding to the key.
//
// ex) mutexs.UnlockAndDelete(key)
func (this *MutexByKey) UnlockAndDelete(key any) {
	if mutex, loaded := this.mutexs.LoadAndDelete(key); loaded {
		mutex.(*Mutex).Unlock()
	} else {
		(&Mutex{}).Unlock()
	}
}

// Delete deletes the mutex corresponding to the key.
//
// ex) mutexs.Delete(key)
func (this *MutexByKey) Delete(key any) {
	this.mutexs.Delete(key)
}
