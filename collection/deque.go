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
func (d *Deque[T]) Front() T {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.datas[0]

}

// Back returns back data.
//
// ex) t := deque.Back()
func (d *Deque[T]) Back() T {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.datas[len(d.datas)-1]
}

// Empty returns whether the queue is empty.
//
// ex) empty := deque.Empty()
func (d *Deque[T]) Empty() bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return len(d.datas) == 0
}

// Size returns the queue size.
//
// ex) size := deque.Size()
func (d *Deque[T]) Size() int {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return len(d.datas)
}

// Clear clears the queue.
//
// ex) deque.Clear()
func (d *Deque[T]) Clear() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.datas = []T{}
}

// PushFront inserts data into the front.
//
// ex) deque.PushFront(1)
func (d *Deque[T]) PushFront(data T) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.datas = append([]T{data}, d.datas...)
}

// PopFront removes front data.
//
// ex) deque.PopFront()
func (d *Deque[T]) PopFront() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if len(d.datas) == 0 {
		return
	}

	d.datas = d.datas[1:]
}

// PushBack inserts data into the back.
//
// ex) deque.PushBack(1)
func (d *Deque[T]) PushBack(data T) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.datas = append(d.datas, data)
}

// PopBack removes back data.
//
// ex) deque.PopBack()
func (d *Deque[T]) PopBack() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if len(d.datas) == 0 {
		return
	}

	d.datas = d.datas[0 : len(d.datas)-1]
}
