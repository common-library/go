# Collection

Thread-safe generic data structures for Go.

## Features

### Queue
- FIFO (First In First Out) data structure
- Thread-safe operations with mutex synchronization
- Generic type support
- Operations: Push, Pop, Front, Back, Size, Empty, Clear

### Deque
- Double-ended queue (deque) data structure
- Thread-safe operations with mutex synchronization
- Generic type support
- Insert/remove from both ends
- Operations: PushFront, PopFront, PushBack, PopBack, Front, Back, Size, Empty, Clear

## Installation

```bash
go get -u github.com/common-library/go/collection
```

## Usage

### Queue

```go
import "github.com/common-library/go/collection"

// Create a queue for integers
var queue collection.Queue[int]

// Push elements
queue.Push(1)
queue.Push(2)
queue.Push(3)

// Check size
size := queue.Size()  // 3

// Access elements
front := queue.Front()  // 1
back := queue.Back()    // 3

// Remove element
queue.Pop()  // Removes 1

// Check if empty
if queue.Empty() {
    fmt.Println("Queue is empty")
}

// Clear all elements
queue.Clear()
```

**Key Functions:**
- `Push(data T)` - Add element to the back
- `Pop()` - Remove element from the front
- `Front() T` - Get front element without removing
- `Back() T` - Get back element without removing
- `Size() int` - Get number of elements
- `Empty() bool` - Check if queue is empty
- `Clear()` - Remove all elements

### Deque

```go
import "github.com/common-library/go/collection"

// Create a deque for strings
var deque collection.Deque[string]

// Push to both ends
deque.PushFront("front1")
deque.PushBack("back1")
deque.PushFront("front2")
deque.PushBack("back2")

// Deque now: [front2, front1, back1, back2]

// Access elements
front := deque.Front()  // "front2"
back := deque.Back()    // "back2"

// Remove from both ends
deque.PopFront()  // Removes "front2"
deque.PopBack()   // Removes "back2"

// Check size
size := deque.Size()  // 2

// Clear all elements
deque.Clear()
```

**Key Functions:**
- `PushFront(data T)` - Add element to the front
- `PopFront()` - Remove element from the front
- `PushBack(data T)` - Add element to the back
- `PopBack()` - Remove element from the back
- `Front() T` - Get front element without removing
- `Back() T` - Get back element without removing
- `Size() int` - Get number of elements
- `Empty() bool` - Check if deque is empty
- `Clear()` - Remove all elements

## Key Differences

| Feature | Queue | Deque |
|---------|-------|-------|
| Insert | Back only | Both ends |
| Remove | Front only | Both ends |
| Access | Front, Back | Front, Back |
| Order | FIFO | Flexible (FIFO/LIFO) |
| Use Cases | Task processing, BFS | Undo/redo, sliding window |

## Implementation Details

### Thread Safety

Both Queue and Deque are thread-safe:
- All operations are protected by mutex locks
- Safe for concurrent access from multiple goroutines
- Uses `github.com/common-library/go/lock` package for synchronization

### Generic Support

Both data structures use Go generics (Go 1.18+):
```go
// Works with any type
var intQueue collection.Queue[int]
var strQueue collection.Queue[string]
var customQueue collection.Queue[MyStruct]

var intDeque collection.Deque[int]
var ptrDeque collection.Deque[*MyStruct]
```

### Performance Characteristics

**Queue:**
- Push: O(1) amortized
- Pop: O(n) - requires slice reslicing
- Front/Back: O(1)
- Size/Empty: O(1)

**Deque:**
- PushFront: O(n) - requires prepending
- PopFront: O(n) - requires slice reslicing
- PushBack: O(1) amortized
- PopBack: O(1)
- Front/Back: O(1)
- Size/Empty: O(1)

## Error Handling

### Panic Prevention

Front() and Back() methods will panic if called on an empty collection:

```go
var queue collection.Queue[int]

// Check before accessing
if !queue.Empty() {
    front := queue.Front()  // Safe
} else {
    fmt.Println("Queue is empty")
}
```

### Safe Pop Operations

Pop operations are safe on empty collections:

```go
var deque collection.Deque[string]

deque.Pop()  // No panic, does nothing
```

## Examples

### BFS (Breadth-First Search) with Queue

```go
type Node struct {
    Value    int
    Children []*Node
}

func BFS(root *Node) []int {
    var queue collection.Queue[*Node]
    var result []int
    
    queue.Push(root)
    
    for !queue.Empty() {
        node := queue.Front()
        queue.Pop()
        
        result = append(result, node.Value)
        
        for _, child := range node.Children {
            queue.Push(child)
        }
    }
    
    return result
}
```

### Sliding Window with Deque

```go
func maxSlidingWindow(nums []int, k int) []int {
    var deque collection.Deque[int]  // Stores indices
    var result []int
    
    for i := 0; i < len(nums); i++ {
        // Remove indices outside window
        for !deque.Empty() && deque.Front() <= i-k {
            deque.PopFront()
        }
        
        // Remove smaller elements
        for !deque.Empty() && nums[deque.Back()] < nums[i] {
            deque.PopBack()
        }
        
        deque.PushBack(i)
        
        if i >= k-1 {
            result = append(result, nums[deque.Front()])
        }
    }
    
    return result
}
```

### Task Queue

```go
type Task struct {
    ID   int
    Name string
}

func ProcessTasks() {
    var taskQueue collection.Queue[Task]
    
    // Producer
    go func() {
        for i := 0; i < 100; i++ {
            taskQueue.Push(Task{ID: i, Name: fmt.Sprintf("Task-%d", i)})
            time.Sleep(10 * time.Millisecond)
        }
    }()
    
    // Consumer
    for {
        if !taskQueue.Empty() {
            task := taskQueue.Front()
            taskQueue.Pop()
            
            // Process task
            fmt.Printf("Processing: %s\n", task.Name)
        }
        
        time.Sleep(50 * time.Millisecond)
    }
}
```

### Undo/Redo Stack with Deque

```go
type Command struct {
    Action string
    Data   interface{}
}

type Editor struct {
    history collection.Deque[Command]
    current int
}

func (e *Editor) Execute(cmd Command) {
    e.history.PushBack(cmd)
    e.current = e.history.Size() - 1
}

func (e *Editor) Undo() {
    if e.current >= 0 {
        e.current--
    }
}

func (e *Editor) Redo() {
    if e.current < e.history.Size()-1 {
        e.current++
    }
}
```

## Best Practices

1. **Check Before Access**
   ```go
   if !queue.Empty() {
       front := queue.Front()
   }
   ```

2. **Clear When Done**
   ```go
   defer queue.Clear()  // Clean up resources
   ```

3. **Use Appropriate Data Structure**
   ```go
   // FIFO only → Queue
   var taskQueue collection.Queue[Task]
   
   // Need both ends → Deque
   var buffer collection.Deque[byte]
   ```

4. **Thread Safety is Built-in**
   ```go
   // Safe without external locking
   go queue.Push(1)
   go queue.Push(2)
   go queue.Pop()
   ```

5. **Generic Type Safety**
   ```go
   // Type-safe at compile time
   var intQueue collection.Queue[int]
   intQueue.Push(42)      // OK
   // intQueue.Push("str") // Compile error
   ```

## Performance Considerations

### When to Use Queue
- Sequential task processing
- BFS algorithms
- Producer-consumer patterns
- Event handling

### When to Use Deque
- Sliding window algorithms
- Undo/redo functionality
- Need efficient insertion/deletion at both ends
- Palindrome checking

### Optimization Tips

For high-performance scenarios:
```go
// Pre-allocate if size is known
// Note: Current implementation doesn't support capacity
// Consider implementing a custom version with pre-allocation

// Batch operations
for _, item := range items {
    queue.Push(item)
}
```

## Dependencies

- `github.com/common-library/go/lock` - Mutex synchronization

## Limitations

1. **No Capacity Limit**: Grows unbounded, may cause memory issues
2. **No Iterator**: No built-in way to iterate over elements
3. **Pop Performance**: O(n) due to slice operations (consider using circular buffer for production)
4. **No Index Access**: Cannot access elements by index

## Alternatives

For production use cases requiring high performance:
- `container/list` - Standard library doubly linked list
- `github.com/gammazero/deque` - High-performance deque
- Custom circular buffer implementation

## Further Reading

- [Go Generics Tutorial](https://go.dev/doc/tutorial/generics)
- [Data Structures in Go](https://golang.org/pkg/container/)
- [Lock package documentation](../lock/)
