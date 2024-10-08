package collection_test

import (
	"testing"

	"github.com/common-library/go/collection"
)

func TestFrontOfQueue(t *testing.T) {
	queue := collection.Queue[int]{}

	queue.Push(1)
	queue.Push(2)

	if queue.Front() != 1 {
		t.Fatal(queue.Front())
	}
}

func TestBackOfQueue(t *testing.T) {
	queue := collection.Queue[int]{}

	queue.Push(1)
	queue.Push(2)

	if queue.Back() != 2 {
		t.Fatal(queue.Back())
	}
}

func TestEmptyOfQueue(t *testing.T) {
	queue := collection.Queue[int]{}

	if queue.Empty() == false {
		t.Fatal(queue.Empty())
	}

	queue.Push(1)

	if queue.Empty() {
		t.Fatal(queue.Empty())
	}
}

func TestSizeOfQueue(t *testing.T) {
	queue := collection.Queue[int]{}

	queue.Push(1)
	if queue.Size() != 1 {
		t.Fatal(queue.Size())
	}

	queue.Push(2)
	if queue.Size() != 2 {
		t.Fatal(queue.Size())
	}
}

func TestClearOfQueue(t *testing.T) {
	queue := collection.Queue[int]{}

	queue.Push(1)
	queue.Push(2)
	if queue.Size() != 2 {
		t.Fatal(queue.Size())
	}

	queue.Clear()
	if queue.Empty() == false {
		t.Fatal(queue.Empty())
	}
}

func TestPushOfQueue(t *testing.T) {
	queue := collection.Queue[int]{}

	queue.Push(1)
	queue.Push(2)
	queue.Push(3)

	if queue.Front() != 1 {
		t.Fatal(queue.Front())
	} else if queue.Back() != 3 {
		t.Fatal(queue.Back())
	}
}

func TestPopOfQueue(t *testing.T) {
	queue := collection.Queue[int]{}

	queue.Pop()

	queue.Push(1)
	if queue.Size() != 1 {
		t.Fatal(queue.Size())
	}

	queue.Pop()
	if queue.Empty() == false {
		t.Fatal(queue.Empty())
	}
}
