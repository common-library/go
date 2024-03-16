// Package collection provides data structure related implementations.
package collection

import "github.com/heaven-chp/common-library-go/lock"

// Queue is struct that provides queue related methods.
type Queue[T any] struct {
	mutex lock.Mutex
	datas []T
}

// Front returns front data.
//
// ex) t := queue.Front()
func (this *Queue[T]) Front() T {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	return this.datas[0]
}

// Back returns back data.
//
// ex) t := queue.Back()
func (this *Queue[T]) Back() T {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	return this.datas[len(this.datas)-1]
}

// Empty returns whether the queue is empty.
//
// ex) empty := queue.Empty()
func (this *Queue[T]) Empty() bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	return len(this.datas) == 0
}

// Size returns the queue size.
//
// ex) size := queue.Size()
func (this *Queue[T]) Size() int {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	return len(this.datas)
}

// Clear clears the queue.
//
// ex) queue.Clear()
func (this *Queue[T]) Clear() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.datas = []T{}
}

// Push inserts data.
//
// ex) queue.Push(1)
func (this *Queue[T]) Push(data T) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.datas = append(this.datas, data)
}

// Pop removes front data.
//
// ex) queue.Pop()
func (this *Queue[T]) Pop() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if len(this.datas) == 0 {
		return
	}

	this.datas = this.datas[1:]
}
