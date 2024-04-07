package collection_test

import (
	"testing"

	"github.com/common-library/go/collection"
)

func TestFrontOfDeque(t *testing.T) {
	deque := collection.Deque[int]{}

	deque.PushFront(1)
	deque.PushBack(2)

	if deque.Front() != 1 {
		t.Fatal("invalid -", deque.Front())
	}
}

func TestBackOfDeque(t *testing.T) {
	deque := collection.Deque[int]{}

	deque.PushFront(1)
	deque.PushBack(2)

	if deque.Back() != 2 {
		t.Fatal("invalid -", deque.Back())
	}
}

func TestEmptyOfDeque(t *testing.T) {
	deque := collection.Deque[int]{}

	if deque.Empty() == false {
		t.Fatal("invalid -", deque.Empty())
	}

	deque.PushFront(1)

	if deque.Empty() {
		t.Fatal("invalid -", deque.Empty())
	}
}

func TestSizeOfDeque(t *testing.T) {
	deque := collection.Deque[int]{}

	deque.PushFront(1)
	if deque.Size() != 1 {
		t.Fatal("invalid -", deque.Size())
	}

	deque.PushFront(2)
	if deque.Size() != 2 {
		t.Fatal("invalid -", deque.Size())
	}
}

func TestClearOfDeque(t *testing.T) {
	deque := collection.Deque[int]{}

	deque.PushFront(1)
	deque.PushFront(2)
	if deque.Size() != 2 {
		t.Fatal("invalid -", deque.Size())
	}

	deque.Clear()
	if deque.Empty() == false {
		t.Fatal("invalid -", deque.Empty())
	}
}

func TestPushFrontOfDeque(t *testing.T) {
	deque := collection.Deque[int]{}

	deque.PushFront(1)
	deque.PushFront(2)
	deque.PushFront(3)

	if deque.Front() != 3 {
		t.Fatal("invalid -", deque.Front())
	} else if deque.Back() != 1 {
		t.Fatal("invalid -", deque.Back())
	}
}

func TestPopFrontOfDeque(t *testing.T) {
	deque := collection.Deque[int]{}

	deque.PopFront()

	deque.PushFront(1)
	if deque.Size() != 1 {
		t.Fatal("invalid -", deque.Size())
	}
	deque.PopFront()
	if deque.Empty() == false {
		t.Fatal("invalid -", deque.Empty())
	}

	deque.PushFront(1)
	deque.PushFront(2)
	deque.PushFront(3)
	if deque.Size() != 3 {
		t.Fatal("invalid -", deque.Size())
	}
	deque.PopFront()
	if deque.Size() != 2 {
		t.Fatal("invalid -", deque.Size())
	} else if deque.Front() != 2 {
		t.Fatal("invalid -", deque.Front())
	} else if deque.Back() != 1 {
		t.Fatal("invalid -", deque.Back())
	}
}

func TestPushBackOfDeque(t *testing.T) {
	deque := collection.Deque[int]{}

	deque.PushBack(1)
	deque.PushBack(2)
	deque.PushBack(3)

	if deque.Front() != 1 {
		t.Fatal("invalid -", deque.Front())
	} else if deque.Back() != 3 {
		t.Fatal("invalid -", deque.Back())
	}
}

func TestPopBackOfDeque(t *testing.T) {
	deque := collection.Deque[int]{}

	deque.PopBack()

	deque.PushBack(1)
	if deque.Size() != 1 {
		t.Fatal("invalid -", deque.Size())
	}
	deque.PopBack()
	if deque.Empty() == false {
		t.Fatal("invalid -", deque.Empty())
	}

	deque.PushBack(1)
	deque.PushBack(2)
	deque.PushBack(3)
	if deque.Size() != 3 {
		t.Fatal("invalid -", deque.Size())
	}
	deque.PopBack()
	if deque.Size() != 2 {
		t.Fatal("invalid -", deque.Size())
	} else if deque.Front() != 1 {
		t.Fatal("invalid -", deque.Front())
	} else if deque.Back() != 2 {
		t.Fatal("invalid -", deque.Back())
	}
}
