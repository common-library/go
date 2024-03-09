package data_structure_test

import (
	"testing"

	data_structure "github.com/heaven-chp/common-library-go/data-structure"
)

func TestFrontOfQueue(t *testing.T) {
	queue := data_structure.Queue[int]{}

	queue.Push(1)
	queue.Push(2)

	if queue.Front() != 1 {
		t.Fatal("invalid -", queue.Front())
	}
}

func TestBackOfQueue(t *testing.T) {
	queue := data_structure.Queue[int]{}

	queue.Push(1)
	queue.Push(2)

	if queue.Back() != 2 {
		t.Fatal("invalid -", queue.Back())
	}
}

func TestEmptyOfQueue(t *testing.T) {
	queue := data_structure.Queue[int]{}

	if queue.Empty() == false {
		t.Fatal("invalid -", queue.Empty())
	}

	queue.Push(1)

	if queue.Empty() {
		t.Fatal("invalid -", queue.Empty())
	}
}

func TestSizeOfQueue(t *testing.T) {
	queue := data_structure.Queue[int]{}

	queue.Push(1)
	if queue.Size() != 1 {
		t.Fatal("invalid -", queue.Size())
	}

	queue.Push(2)
	if queue.Size() != 2 {
		t.Fatal("invalid -", queue.Size())
	}
}

func TestClearOfQueue(t *testing.T) {
	queue := data_structure.Queue[int]{}

	queue.Push(1)
	queue.Push(2)
	if queue.Size() != 2 {
		t.Fatal("invalid -", queue.Size())
	}

	queue.Clear()
	if queue.Empty() == false {
		t.Fatal("invalid -", queue.Empty())
	}
}

func TestPushOfQueue(t *testing.T) {
	queue := data_structure.Queue[int]{}

	queue.Push(1)
	queue.Push(2)
	queue.Push(3)

	if queue.Front() != 1 {
		t.Fatal("invalid -", queue.Front())
	} else if queue.Back() != 3 {
		t.Fatal("invalid -", queue.Back())
	}
}

func TestPopOfQueue(t *testing.T) {
	queue := data_structure.Queue[int]{}

	queue.Pop()

	queue.Push(1)
	if queue.Size() != 1 {
		t.Fatal("invalid -", queue.Size())
	}

	queue.Pop()
	if queue.Empty() == false {
		t.Fatal("invalid -", queue.Empty())
	}
}
