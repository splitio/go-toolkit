package queuecache

import (
	"errors"
	"fmt"
	"math"
	"sync"
)

// RefillError struct to be returned when the refill function panics
type RefillError struct {
	OriginalPanic interface{}
}

func (e *RefillError) Error() string {
	return "Supplied refilling function panicked. See `.OriginalPanic` property to get panicked content"
}

// MessagesDroppedError is the Error to be returned when messages fail to be added to the queue.
type MessagesDroppedError struct {
	MessagesDropped int
}

func (e *MessagesDroppedError) Error() string {
	return fmt.Sprintf(
		"%d messages were dropped. Please report this error as it's most likely a bug in the library",
		e.MessagesDropped,
	)
}

// InMemoryQueueCacheOverlay offers an in-memory queue that gets re-populated whenever it runs out of items
type InMemoryQueueCacheOverlay[T any] struct {
	maxSize      int
	writeCursor  int
	readCursor   int
	queue        []T
	lock         sync.Mutex
	refillCustom func(count int) ([]T, error)
}

// New creates a new InMemoryQueueCacheOverlay
func New[T any](maxSize int, refillFunc func(count int) ([]T, error)) *InMemoryQueueCacheOverlay[T] {
	return &InMemoryQueueCacheOverlay[T]{
		queue:        make([]T, maxSize),
		maxSize:      maxSize,
		writeCursor:  0,
		readCursor:   0,
		refillCustom: refillFunc,
	}
}

// Count returns the number of cached items
func (i *InMemoryQueueCacheOverlay[T]) Count() int {
	if i.writeCursor == i.readCursor {
		return 0
	} else if i.writeCursor > i.readCursor {
		return i.writeCursor - i.readCursor
	}
	return i.maxSize - (i.readCursor - i.writeCursor)
}

func (i *InMemoryQueueCacheOverlay[T]) write(elem T) error {
	if ((i.writeCursor + 1) % i.maxSize) == i.readCursor {
		return errors.New("QUEUE_FULL")
	}

	i.queue[i.writeCursor] = elem
	i.writeCursor = (i.writeCursor + 1) % i.maxSize
	return nil
}

func (i *InMemoryQueueCacheOverlay[T]) read() (T, error) {
	if i.readCursor == i.writeCursor {
        var t T
		return t, errors.New("QUEUE_EMPTY")
	}

	toReturn := i.queue[i.readCursor]
	i.readCursor = (i.readCursor + 1) % i.maxSize
	return toReturn, nil
}

func (i *InMemoryQueueCacheOverlay[T]) refillWrapper(count int) (result []T, err error) {
	defer func() {
		if r := recover(); r != nil {
			result = nil
			err = &RefillError{OriginalPanic: r}
		}
	}()

	return i.refillCustom(count)

}

// Fetch items (will re-populate if necessary)
func (i *InMemoryQueueCacheOverlay[T]) Fetch(requestedCount int) ([]T, error) {
	defer i.lock.Unlock()
	i.lock.Lock()

	dropped := 0
	if i.Count() < requestedCount {
		toAdd, err := i.refillWrapper(i.maxSize - i.Count() - 1)
		if err != nil {
			return nil, err
		}
		for _, item := range toAdd {
			err = i.write(item)
			if err != nil {
				dropped++
			}
		}
	}

	toReturn := make([]T, int(math.Min(float64(requestedCount), float64(i.Count()))))
	for index := 0; index < len(toReturn); index++ {
		elem, err := i.read()
		if err != nil {
			return toReturn[0:index], nil
		}
		toReturn[index] = elem
	}

	if dropped > 0 {
		return toReturn, &MessagesDroppedError{MessagesDropped: dropped}
	}
	return toReturn, nil

}
