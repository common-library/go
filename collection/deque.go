// Package collection provides data structure related implementations.
package collection

import "github.com/common-library/go/lock"

// Deque is struct that provides deque related methods.
type Deque[T any] struct {
	mutex lock.Mutex
	datas []T
}

// Front returns front data.
//
// ex) t := deque.Front()
func (this *Deque[T]) Front() T {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	return this.datas[0]

}

// Back returns back data.
//
// ex) t := deque.Back()
func (this *Deque[T]) Back() T {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	return this.datas[len(this.datas)-1]
}

// Empty returns whether the queue is empty.
//
// ex) empty := deque.Empty()
func (this *Deque[T]) Empty() bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	return len(this.datas) == 0
}

// Size returns the queue size.
//
// ex) size := deque.Size()
func (this *Deque[T]) Size() int {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	return len(this.datas)
}

// Clear clears the queue.
//
// ex) deque.Clear()
func (this *Deque[T]) Clear() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.datas = []T{}
}

// PushFront inserts data into the front.
//
// ex) deque.PushFront(1)
func (this *Deque[T]) PushFront(data T) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.datas = append([]T{data}, this.datas...)
}

// PopFront removes front data.
//
// ex) deque.PopFront()
func (this *Deque[T]) PopFront() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if len(this.datas) == 0 {
		return
	}

	this.datas = this.datas[1:]
}

// PushBack inserts data into the back.
//
// ex) deque.PushBack(1)
func (this *Deque[T]) PushBack(data T) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.datas = append(this.datas, data)
}

// PopBack removes back data.
//
// ex) deque.PopBack()
func (this *Deque[T]) PopBack() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if len(this.datas) == 0 {
		return
	}

	this.datas = this.datas[0 : len(this.datas)-1]
}
