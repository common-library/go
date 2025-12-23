// Package collection provides thread-safe data structure implementations.
//
// This package offers generic data structures with built-in synchronization
// using the common-library lock package.
//
// Features:
//   - Thread-safe Queue (FIFO) with generics
//   - Thread-safe Deque (double-ended queue) with generics
//   - Automatic mutex-based synchronization
//   - Type-safe operations using Go generics
//
// Example usage:
//
//	var d collection.Deque[string]
//	d.PushBack("hello")
//	back := d.Back()
//	d.PopBack()
package collection

import "github.com/common-library/go/lock"

// Deque is struct that provides deque related methods.
type Deque[T any] struct {
	mutex lock.Mutex
	datas []T
}

// Front returns the front element of the deque without removing it.
//
// This is a thread-safe operation. The caller should ensure the deque is not empty
// before calling this method to avoid index out of range panics.
//
// Returns the element at the front of the deque.
//
// Example:
//
//	front := deque.Front()
func (d *Deque[T]) Front() T {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.datas[0]

}

// Back returns the back element of the deque without removing it.
//
// This is a thread-safe operation. The caller should ensure the deque is not empty
// before calling this method to avoid index out of range panics.
//
// Returns the element at the back of the deque.
//
// Example:
//
//	back := deque.Back()
func (d *Deque[T]) Back() T {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.datas[len(d.datas)-1]
}

// Empty returns true if the deque contains no elements.
//
// This is a thread-safe operation.
//
// Returns true if the deque is empty, false otherwise.
//
// Example:
//
//	if deque.Empty() {
//	    fmt.Println("Deque is empty")
//	}
func (d *Deque[T]) Empty() bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return len(d.datas) == 0
}

// Size returns the number of elements in the deque.
//
// This is a thread-safe operation.
//
// Returns the current size of the deque.
//
// Example:
//
//	size := deque.Size()
//	fmt.Printf("Deque has %d elements\n", size)
func (d *Deque[T]) Size() int {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return len(d.datas)
}

// Clear removes all elements from the deque.
//
// This is a thread-safe operation. After calling Clear, the deque will be empty
// and Size will return 0.
//
// Example:
//
//	deque.Clear()
func (d *Deque[T]) Clear() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.datas = []T{}
}

// PushFront inserts an element at the front of the deque.
//
// This is a thread-safe operation. Elements can be added to either end of the deque.
//
// Parameters:
//   - data: the element to add to the front of the deque
//
// Example:
//
//	deque.PushFront(42)
//	deque.PushFront("hello")
func (d *Deque[T]) PushFront(data T) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.datas = append([]T{data}, d.datas...)
}

// PopFront removes the element at the front of the deque.
//
// This is a thread-safe operation. If the deque is empty, this method does nothing.
// Elements can be removed from either end of the deque.
//
// Example:
//
//	deque.PopFront()
func (d *Deque[T]) PopFront() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if len(d.datas) == 0 {
		return
	}

	d.datas = d.datas[1:]
}

// PushBack inserts an element at the back of the deque.
//
// This is a thread-safe operation. Elements can be added to either end of the deque.
//
// Parameters:
//   - data: the element to add to the back of the deque
//
// Example:
//
//	deque.PushBack(42)
//	deque.PushBack("hello")
func (d *Deque[T]) PushBack(data T) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.datas = append(d.datas, data)
}

// PopBack removes the element at the back of the deque.
//
// This is a thread-safe operation. If the deque is empty, this method does nothing.
// Elements can be removed from either end of the deque.
//
// Example:
//
//	deque.PopBack()
func (d *Deque[T]) PopBack() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if len(d.datas) == 0 {
		return
	}

	d.datas = d.datas[0 : len(d.datas)-1]
}
