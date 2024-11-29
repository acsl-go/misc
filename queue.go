package misc

import "sync"

type _queueItem[T comparable] struct {
	data T
	next *_queueItem[T]
}

type Queue[T comparable] struct {
	head *_queueItem[T]
	tail *_queueItem[T]
	Pool *QueueItemPool[T]
}

type QueueItemPool[T comparable] struct {
	pool *sync.Pool
}

func NewQueueItemPool[T comparable]() *QueueItemPool[T] {
	return &QueueItemPool[T]{
		pool: &sync.Pool{
			New: func() interface{} {
				return &_queueItem[T]{}
			},
		},
	}
}

func (q *Queue[T]) _allocItem(data T) *_queueItem[T] {
	if q.Pool == nil {
		return &_queueItem[T]{
			data: data,
		}
	}
	item := q.Pool.pool.Get().(*_queueItem[T])
	item.data = data
	return item
}

func (q *Queue[T]) _freeItem(item *_queueItem[T]) {
	if q.Pool != nil {
		var zero T
		item.data = zero
		q.Pool.pool.Put(item)
	}
}

func (q *Queue[T]) Enqueue(data T) {
	item := q._allocItem(data)
	if q.head == nil {
		q.head = item
		q.tail = item
	} else {
		q.tail.next = item
		q.tail = item
	}
}

func (q *Queue[T]) Dequeue() T {
	var retVal T
	if q.head == nil {
		return retVal
	}
	item := q.head
	q.head = item.next
	if q.head == nil {
		q.tail = nil
	}
	retVal = item.data
	q._freeItem(item)
	return retVal
}

func (q *Queue[T]) Remove(data T) {
	var prev *_queueItem[T]
	for item := q.head; item != nil; item = item.next {
		if item.data == data {
			if prev == nil {
				q.head = item.next
			} else {
				prev.next = item.next
			}
			if item == q.tail {
				q.tail = prev
			}
			item.next = nil
			q._freeItem(item)
			return
		}
		prev = item
	}
}

func (q *Queue[T]) RemoveEx(data T, isEqual func(T, T) bool) {
	var prev *_queueItem[T]
	for item := q.head; item != nil; item = item.next {
		if isEqual(item.data, data) {
			if prev == nil {
				q.head = item.next
			} else {
				prev.next = item.next
			}
			if item == q.tail {
				q.tail = prev
			}
			item.next = nil
			q._freeItem(item)
			return
		}
		prev = item
	}
}

func (q *Queue[T]) Clear() {
	p := q.head
	t := p
	for p != nil {
		t = p
		p = p.next
		q._freeItem(t)
	}
	q.head = nil
	q.tail = nil
}

func (q *Queue[T]) ClearEx(f func(T)) {
	p := q.head
	t := p
	for p != nil {
		t = p
		p = p.next
		f(t.data)
		q._freeItem(t)
	}
	q.head = nil
	q.tail = nil
}

func (q *Queue[T]) IsEmpty() bool {
	return q.head == nil
}

func (q *Queue[T]) First() T {
	var zero T
	if q.head == nil {
		return zero
	}
	return q.head.data
}

func (q *Queue[T]) Last() T {
	var zero T
	if q.tail == nil {
		return zero
	}
	return q.tail.data
}

func (q *Queue[T]) Contains(data T) bool {
	for item := q.head; item != nil; item = item.next {
		if item.data == data {
			return true
		}
	}
	return false
}

func (q *Queue[T]) ContainsEx(data T, isEqual func(T, T) bool) bool {
	for item := q.head; item != nil; item = item.next {
		if isEqual(item.data, data) {
			return true
		}
	}
	return false
}

func (q *Queue[T]) Size() int {
	size := 0
	for item := q.head; item != nil; item = item.next {
		size++
	}
	return size
}

func (q *Queue[T]) ToSlice() []T {
	slice := make([]T, 0, q.Size())
	for item := q.head; item != nil; item = item.next {
		slice = append(slice, item.data)
	}
	return slice
}

func (q *Queue[T]) ForEach(f func(T) bool) {
	for item := q.head; item != nil; item = item.next {
		if !f(item.data) {
			break
		}
	}
}

func (q *Queue[T]) ForEachEx(f func(T, int) bool) {
	index := 0
	for item := q.head; item != nil; item = item.next {
		if !f(item.data, index) {
			break
		}
		index++
	}
}
