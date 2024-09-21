package lock_test

import (
	"sync"
	"testing"

	"github.com/common-library/go/lock"
)

func TestMutex(t *testing.T) {
	t.Parallel()

	mutex := lock.Mutex{}

	mutex.Lock()
	mutex.Unlock()

	if mutex.TryLock() == false {
		t.Fatal("invalid")
	}
	if mutex.TryLock() {
		t.Fatal("invalid")
	}
	mutex.Unlock()

	const count = 1000
	result := 0
	wg := new(sync.WaitGroup)
	for i := 1; i <= count; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			mutex.Lock()
			defer mutex.Unlock()

			result += index
		}(i)
	}
	wg.Wait()

	if result != count*(count+1)/2 {
		t.Fatal(count, result)
	}
}

func TestMutexByKey(t *testing.T) {
	t.Parallel()

	mutexs := lock.MutexByKey{}

	test := func(key any) {
		mutexs.Lock(key)
		mutexs.Unlock(key)

		mutexs.Delete(key)

		if mutexs.TryLock(key) == false {
			t.Fatal("invalid")
		}
		if mutexs.TryLock(key) {
			t.Fatal("invalid")
		}
		mutexs.UnlockAndDelete(key)

		const count = 1000
		result := 0
		wg := new(sync.WaitGroup)
		for i := 1; i <= count; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				mutexs.Lock(key)
				defer mutexs.Unlock(key)

				result += index
			}(i)
		}
		wg.Wait()
		if result != count*(count+1)/2 {
			t.Fatal(count, result)
		}
	}

	test(1)
	test("a")
}
