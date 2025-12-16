// Package collection provides data structure related implementations.
package collection

import "github.com/common-library/go/lock"

// Queue is struct that provides queue related methods.
type Queue[T any] struct {
	mutex lock.Mutex
	datas []T
}

// Front returns front data.
//
// ex) t := queue.Front()
func (q *Queue[T]) Front() T {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return q.datas[0]
}

// Back returns back data.
//
// ex) t := queue.Back()
func (q *Queue[T]) Back() T {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return q.datas[len(q.datas)-1]
}

// Empty returns whether the queue is empty.
//
// ex) empty := queue.Empty()
func (q *Queue[T]) Empty() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return len(q.datas) == 0
}

// Size returns the queue size.
//
// ex) size := queue.Size()
func (q *Queue[T]) Size() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return len(q.datas)
}

// Clear clears the queue.
//
// ex) queue.Clear()
func (q *Queue[T]) Clear() {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.datas = []T{}
}

// Push inserts data.
//
// ex) queue.Push(1)
func (q *Queue[T]) Push(data T) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.datas = append(q.datas, data)
}

// Pop removes front data.
//
// ex) queue.Pop()
func (q *Queue[T]) Pop() {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.datas) == 0 {
		return
	}

	q.datas = q.datas[1:]
}
