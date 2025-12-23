// Package lock provides mutex implementations for synchronization.
//
// This package offers enhanced mutex functionality including basic mutex
// operations and key-based mutex management for concurrent access control.
//
// Features:
//   - Basic mutex wrapper with Lock/TryLock/Unlock
//   - Key-based mutex map for managing multiple locks
//   - Thread-safe mutex storage using sync.Map
//   - TryLock support for non-blocking lock acquisition
//
// Example:
//
//	var mu lock.Mutex
//	mu.Lock()
//	defer mu.Unlock()
//	// Critical section
//
//	var mutexes lock.MutexByKey
//	mutexes.Lock("resource-1")
//	defer mutexes.Unlock("resource-1")
package lock

import "sync"

// Mutex is a wrapper around sync.Mutex providing mutex operations.
//
// This type provides the same functionality as sync.Mutex with additional
// methods like TryLock for non-blocking lock attempts.
//
// Example:
//
//	var mu lock.Mutex
//
//	// Blocking lock
//	mu.Lock()
//	defer mu.Unlock()
//	// Critical section
//
//	// Non-blocking lock
//	if mu.TryLock() {
//	    defer mu.Unlock()
//	    // Critical section
//	} else {
//	    // Lock not acquired
//	}
type Mutex struct {
	mutex sync.Mutex
}

// Lock acquires the mutex, blocking until the lock is available.
//
// If the mutex is already locked by another goroutine, Lock blocks until
// the mutex becomes available. It is a run-time error if m is already
// locked by the calling goroutine.
//
// Behavior:
//   - Blocks the calling goroutine until lock is acquired
//   - Must be paired with Unlock
//   - Calling Lock on an already-locked mutex by the same goroutine causes deadlock
//
// Example:
//
//	var mu lock.Mutex
//
//	// Acquire lock
//	mu.Lock()
//	defer mu.Unlock() // Ensure unlock
//
//	// Critical section - only one goroutine executes this at a time
//	fmt.Println("Processing...")
//	sharedResource++
//
// Example with goroutines:
//
//	var mu lock.Mutex
//	var counter int
//
//	for i := 0; i < 10; i++ {
//	    go func() {
//	        mu.Lock()
//	        counter++
//	        mu.Unlock()
//	    }()
//	}
func (m *Mutex) Lock() {
	m.mutex.Lock()
}

// TryLock attempts to acquire the mutex without blocking.
//
// TryLock tries to lock the mutex and reports whether it succeeded.
// Unlike Lock, TryLock does not block if the mutex is already locked.
//
// Returns:
//   - bool: true if the lock was acquired, false otherwise
//
// Behavior:
//   - Returns immediately (non-blocking)
//   - Returns true if lock acquired successfully
//   - Returns false if mutex is already locked
//   - If true is returned, Unlock must be called to release the lock
//
// Example:
//
//	var mu lock.Mutex
//
//	if mu.TryLock() {
//	    defer mu.Unlock()
//	    fmt.Println("Lock acquired, processing...")
//	    // Critical section
//	} else {
//	    fmt.Println("Lock not available, skipping...")
//	    // Alternative action
//	}
//
// Example with retry:
//
//	var mu lock.Mutex
//
//	for i := 0; i < 3; i++ {
//	    if mu.TryLock() {
//	        defer mu.Unlock()
//	        // Process
//	        break
//	    }
//	    time.Sleep(100 * time.Millisecond)
//	}
//
// Example with timeout pattern:
//
//	var mu lock.Mutex
//
//	timeout := time.After(5 * time.Second)
//	ticker := time.NewTicker(100 * time.Millisecond)
//	defer ticker.Stop()
//
//	for {
//	    select {
//	    case <-timeout:
//	        fmt.Println("Timeout")
//	        return
//	    case <-ticker.C:
//	        if mu.TryLock() {
//	            defer mu.Unlock()
//	            // Process
//	            return
//	        }
//	    }
//	}
func (m *Mutex) TryLock() bool {
	return m.mutex.TryLock()
}

// Unlock releases the mutex.
//
// Unlock unlocks the mutex. It is a run-time error if the mutex is not
// locked on entry to Unlock.
//
// Behavior:
//   - Releases the mutex for other goroutines to acquire
//   - Must be called by the same goroutine that called Lock
//   - Calling Unlock on an unlocked mutex causes panic
//   - Should be paired with Lock or successful TryLock
//
// A locked Mutex is not associated with a particular goroutine.
// It is allowed for one goroutine to lock a Mutex and then arrange
// for another goroutine to unlock it (though this is generally not recommended).
//
// Example with defer:
//
//	var mu lock.Mutex
//
//	mu.Lock()
//	defer mu.Unlock() // Automatically unlocks when function returns
//
//	// Critical section
//	if err := doSomething(); err != nil {
//	    return err // Unlock happens automatically
//	}
//
// Example without defer:
//
//	var mu lock.Mutex
//
//	mu.Lock()
//	// Critical section
//	processData()
//	mu.Unlock()
//
// Example with early unlock:
//
//	var mu lock.Mutex
//
//	mu.Lock()
//	data := readSharedData()
//	mu.Unlock() // Unlock early to reduce contention
//
//	// Process data without holding the lock
//	result := processData(data)
func (m *Mutex) Unlock() {
	m.mutex.Unlock()
}

// MutexByKey manages multiple mutexes indexed by keys.
//
// This type maintains a thread-safe map of mutexes, allowing different
// resources to be locked independently using unique keys. It's useful
// when you need to synchronize access to multiple resources without
// locking all of them at once.
//
// The underlying storage uses sync.Map for thread-safe concurrent access.
// Mutexes are created on-demand when first accessed.
//
// Example:
//
//	var mutexes lock.MutexByKey
//
//	// Lock different resources independently
//	mutexes.Lock("user:123")
//	defer mutexes.Unlock("user:123")
//	// Update user 123
//
//	mutexes.Lock("account:456")
//	defer mutexes.Unlock("account:456")
//	// Update account 456
type MutexByKey struct {
	mutexs sync.Map
}

// Lock acquires the mutex associated with the given key.
//
// If no mutex exists for the key, a new one is created automatically.
// This operation is thread-safe and multiple goroutines can call Lock
// with different keys simultaneously without blocking each other.
//
// Parameters:
//   - key: Any comparable value to identify the mutex (typically string or int)
//
// Behavior:
//   - Creates mutex on-demand if it doesn't exist
//   - Blocks until the mutex for the specific key is acquired
//   - Must be paired with Unlock or UnlockAndDelete using the same key
//   - Thread-safe: multiple keys can be locked concurrently
//
// Example:
//
//	var mutexes lock.MutexByKey
//
//	// Lock specific user
//	mutexes.Lock("user:123")
//	defer mutexes.Unlock("user:123")
//
//	// Update user data
//	updateUser(123)
//
// Example with multiple resources:
//
//	var mutexes lock.MutexByKey
//
//	// Different resources can be locked independently
//	go func() {
//	    mutexes.Lock("resource-A")
//	    defer mutexes.Unlock("resource-A")
//	    // Process A
//	}()
//
//	go func() {
//	    mutexes.Lock("resource-B")
//	    defer mutexes.Unlock("resource-B")
//	    // Process B (runs concurrently with A)
//	}()
//
// Example with struct key:
//
//	type ResourceKey struct {
//	    Type string
//	    ID   int
//	}
//
//	key := ResourceKey{Type: "user", ID: 123}
//	mutexes.Lock(key)
//	defer mutexes.Unlock(key)
func (mbk *MutexByKey) Lock(key any) {
	mutex, _ := mbk.mutexs.LoadOrStore(key, &Mutex{})

	mutex.(*Mutex).Lock()
}

// TryLock attempts to acquire the mutex for the given key without blocking.
//
// TryLock tries to lock the mutex associated with the key and reports
// whether it succeeded. If no mutex exists for the key, a new one is
// created automatically.
//
// Parameters:
//   - key: Any comparable value to identify the mutex
//
// Returns:
//   - bool: true if the lock was acquired, false otherwise
//
// Behavior:
//   - Creates mutex on-demand if it doesn't exist
//   - Returns immediately without blocking
//   - Returns true if lock acquired successfully
//   - Returns false if mutex for the key is already locked
//   - If true is returned, Unlock or UnlockAndDelete must be called
//
// Example:
//
//	var mutexes lock.MutexByKey
//
//	if mutexes.TryLock("user:123") {
//	    defer mutexes.Unlock("user:123")
//	    fmt.Println("Processing user 123")
//	    updateUser(123)
//	} else {
//	    fmt.Println("User 123 is being processed, skip")
//	}
//
// Example with fallback:
//
//	var mutexes lock.MutexByKey
//
//	if mutexes.TryLock("resource-A") {
//	    defer mutexes.Unlock("resource-A")
//	    processResourceA()
//	} else {
//	    // Process alternative resource if A is busy
//	    mutexes.Lock("resource-B")
//	    defer mutexes.Unlock("resource-B")
//	    processResourceB()
//	}
//
// Example with retry:
//
//	var mutexes lock.MutexByKey
//	key := "critical-resource"
//
//	for i := 0; i < 3; i++ {
//	    if mutexes.TryLock(key) {
//	        defer mutexes.Unlock(key)
//	        processCriticalResource()
//	        break
//	    }
//	    time.Sleep(100 * time.Millisecond)
//	}
func (mbk *MutexByKey) TryLock(key any) bool {
	mutex, _ := mbk.mutexs.LoadOrStore(key, &Mutex{})

	return mutex.(*Mutex).TryLock()
}

// Unlock releases the mutex associated with the given key.
//
// Unlock unlocks the mutex for the specified key. The mutex must be locked
// before calling Unlock. If no mutex exists for the key, a new unlocked one
// is created, which may cause a panic.
//
// Parameters:
//   - key: The key identifying the mutex to unlock
//
// Behavior:
//   - Releases the mutex for other goroutines to acquire
//   - Must be called with the same key used for Lock/TryLock
//   - Calling Unlock on an unlocked mutex causes panic
//   - The mutex remains in the map after unlocking
//
// Note: Use UnlockAndDelete if you want to remove the mutex after unlocking
// to free memory for mutexes that won't be used again.
//
// Example with defer:
//
//	var mutexes lock.MutexByKey
//
//	mutexes.Lock("user:123")
//	defer mutexes.Unlock("user:123")
//
//	// Update user
//	updateUser(123)
//
// Example without defer:
//
//	var mutexes lock.MutexByKey
//
//	mutexes.Lock("account:456")
//	processAccount(456)
//	mutexes.Unlock("account:456")
//
// Example with early unlock:
//
//	var mutexes lock.MutexByKey
//
//	mutexes.Lock("data:789")
//	data := readData(789)
//	mutexes.Unlock("data:789") // Unlock early
//
//	// Process without holding the lock
//	result := processData(data)
func (mbk *MutexByKey) Unlock(key any) {
	mutex, _ := mbk.mutexs.LoadOrStore(key, &Mutex{})

	mutex.(*Mutex).Unlock()
}

// UnlockAndDelete releases and removes the mutex for the given key.
//
// This method unlocks the mutex and then deletes it from the internal map,
// freeing memory. This is useful for one-time locks or when you know a
// particular key won't be used again.
//
// Parameters:
//   - key: The key identifying the mutex to unlock and delete
//
// Behavior:
//   - Unlocks the mutex if it exists
//   - Removes the mutex from the internal map
//   - Safe to call even if the mutex doesn't exist
//   - Calling on an unlocked mutex causes panic (before deletion)
//
// Use this instead of Unlock when:
//   - Processing one-time resources
//   - Managing many temporary locks
//   - Memory efficiency is important
//
// Example for temporary locks:
//
//	var mutexes lock.MutexByKey
//
//	requestID := generateRequestID()
//	mutexes.Lock(requestID)
//	defer mutexes.UnlockAndDelete(requestID) // Clean up after request
//
//	processRequest(requestID)
//
// Example for batch processing:
//
//	var mutexes lock.MutexByKey
//
//	for _, item := range items {
//	    key := fmt.Sprintf("item:%d", item.ID)
//	    mutexes.Lock(key)
//
//	    processItem(item)
//
//	    // Clean up since we won't need this lock again
//	    mutexes.UnlockAndDelete(key)
//	}
//
// Example for session management:
//
//	var mutexes lock.MutexByKey
//
//	func handleSession(sessionID string) {
//	    mutexes.Lock(sessionID)
//	    defer mutexes.UnlockAndDelete(sessionID)
//
//	    // Process session
//	    // Lock is cleaned up when session ends
//	}
//
// Comparison with Unlock:
//
//	// Using Unlock - mutex stays in memory
//	mutexes.Lock("temp-key")
//	mutexes.Unlock("temp-key") // Mutex still in map
//
//	// Using UnlockAndDelete - mutex is removed
//	mutexes.Lock("temp-key")
//	mutexes.UnlockAndDelete("temp-key") // Mutex removed from map
func (mbk *MutexByKey) UnlockAndDelete(key any) {
	if mutex, loaded := mbk.mutexs.LoadAndDelete(key); loaded {
		mutex.(*Mutex).Unlock()
	} else {
		(&Mutex{}).Unlock()
	}
}

// Delete removes the mutex for the given key from the internal map.
//
// This method removes the mutex without unlocking it first. Use with caution,
// as deleting a locked mutex can lead to resource leaks if other goroutines
// are waiting for it.
//
// Parameters:
//   - key: The key identifying the mutex to delete
//
// Behavior:
//   - Removes the mutex from the internal map
//   - Does NOT unlock the mutex before deletion
//   - Safe to call even if the key doesn't exist
//   - No effect if the key is not in the map
//
// Warning: Only use Delete when you're certain the mutex is not locked
// or when cleaning up abandoned locks. Prefer UnlockAndDelete for normal
// cleanup after releasing a lock.
//
// Example for cleanup:
//
//	var mutexes lock.MutexByKey
//
//	// After determining a lock is no longer needed
//	// (and is definitely not locked)
//	mutexes.Delete("old-key")
//
// Example for resetting state:
//
//	var mutexes lock.MutexByKey
//
//	func reset() {
//	    // Clear all mutexes
//	    for _, key := range getAllKeys() {
//	        mutexes.Delete(key)
//	    }
//	}
//
// Example for conditional cleanup:
//
//	var mutexes lock.MutexByKey
//
//	if !isActive(key) {
//	    // Remove inactive resource's mutex
//	    mutexes.Delete(key)
//	}
//
// Preferred pattern (use UnlockAndDelete instead):
//
//	// Good: Unlock and delete together
//	mutexes.UnlockAndDelete(key)
//
//	// Risky: Delete without unlock
//	// Only do this if you're absolutely sure it's not locked
//	mutexes.Delete(key)
func (mbk *MutexByKey) Delete(key any) {
	mbk.mutexs.Delete(key)
}
