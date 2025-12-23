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
//	var q collection.Queue[int]
//	q.Push(1)
//	front := q.Front()
//	q.Pop()
package collection

import "github.com/common-library/go/lock"

// Queue is struct that provides queue related methods.
type Queue[T any] struct {
	mutex lock.Mutex
	datas []T
}

// Front returns the front element of the queue without removing it.
//
// This is a thread-safe operation. The caller should ensure the queue is not empty
// before calling this method to avoid index out of range panics.
//
// Returns the element at the front of the queue.
//
// Example:
//
//	front := queue.Front()
func (q *Queue[T]) Front() T {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return q.datas[0]
}

// Back returns the back element of the queue without removing it.
//
// This is a thread-safe operation. The caller should ensure the queue is not empty
// before calling this method to avoid index out of range panics.
//
// Returns the element at the back of the queue.
//
// Example:
//
//	back := queue.Back()
func (q *Queue[T]) Back() T {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return q.datas[len(q.datas)-1]
}

// Empty returns true if the queue contains no elements.
//
// This is a thread-safe operation.
//
// Returns true if the queue is empty, false otherwise.
//
// Example:
//
//	if queue.Empty() {
//	    fmt.Println("Queue is empty")
//	}
func (q *Queue[T]) Empty() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return len(q.datas) == 0
}

// Size returns the number of elements in the queue.
//
// This is a thread-safe operation.
//
// Returns the current size of the queue.
//
// Example:
//
//	size := queue.Size()
//	fmt.Printf("Queue has %d elements\n", size)
func (q *Queue[T]) Size() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return len(q.datas)
}

// Clear removes all elements from the queue.
//
// This is a thread-safe operation. After calling Clear, the queue will be empty
// and Size will return 0.
//
// Example:
//
//	queue.Clear()
func (q *Queue[T]) Clear() {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.datas = []T{}
}

// Push inserts an element at the back of the queue.
//
// This is a thread-safe operation. Elements are added to the back and removed
// from the front, implementing FIFO (First In First Out) behavior.
//
// Parameters:
//   - data: the element to add to the queue
//
// Example:
//
//	queue.Push(42)
//	queue.Push("hello")
func (q *Queue[T]) Push(data T) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.datas = append(q.datas, data)
}

// Pop removes the element at the front of the queue.
//
// This is a thread-safe operation. If the queue is empty, this method does nothing.
// Elements are removed from the front, implementing FIFO (First In First Out) behavior.
//
// Example:
//
//	queue.Pop()
func (q *Queue[T]) Pop() {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.datas) == 0 {
		return
	}

	q.datas = q.datas[1:]
}
