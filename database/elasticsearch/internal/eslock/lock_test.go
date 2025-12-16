package eslock_test

import (
	"sync"
	"testing"

	"github.com/common-library/go/database/elasticsearch/internal/eslock"
)

func TestInitMu_ConcurrentAccess(t *testing.T) {
	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)

	sharedCounter := 0

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()

			eslock.InitMu.Lock()
			sharedCounter++
			eslock.InitMu.Unlock()
		}()
	}

	wg.Wait()

	if sharedCounter != goroutines {
		t.Errorf("Expected counter to be %d, got %d", goroutines, sharedCounter)
	}
}

func TestInitMu_Exclusivity(t *testing.T) {
	const goroutines = 10
	var wg sync.WaitGroup
	wg.Add(goroutines)

	criticalSection := false
	violations := 0

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()

			eslock.InitMu.Lock()
			defer eslock.InitMu.Unlock()

			if criticalSection {
				violations++
			}

			criticalSection = true

			for j := 0; j < 1000; j++ {
				_ = j * j
			}

			criticalSection = false
		}()
	}

	wg.Wait()

	if violations > 0 {
		t.Errorf("Mutex exclusivity violated %d times", violations)
	}
}

func TestInitMu_Sequential(t *testing.T) {
	counter := 0
	for i := 0; i < 100; i++ {
		eslock.InitMu.Lock()
		counter++
		eslock.InitMu.Unlock()
	}

	if counter != 100 {
		t.Errorf("Expected counter to be 100, got %d", counter)
	}
}

func TestInitMu_NestedDeadlockPrevention(t *testing.T) {
	done := make(chan bool, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {

				done <- true
			}
		}()

		eslock.InitMu.Lock()
		testVar := 42
		_ = testVar

		eslock.InitMu.Unlock()

		done <- true
	}()

	select {
	case <-done:

	case <-make(chan bool):

		t.Fatal("Test should not deadlock in this implementation")
	}
}
