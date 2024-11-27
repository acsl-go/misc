package misc

type _queueItem[T comparable] struct {
	data T
	next *_queueItem[T]
}

type Queue[T comparable] struct {
	head *_queueItem[T]
	tail *_queueItem[T]
}

func (q *Queue[T]) Enqueue(data T) {
	item := &_queueItem[T]{data: data}
	if q.head == nil {
		q.head = item
		q.tail = item
	} else {
		q.tail.next = item
		q.tail = item
	}
}

func (q *Queue[T]) Dequeue() T {
	var zero T
	if q.head == nil {
		return zero
	}
	item := q.head
	q.head = item.next
	if q.head == nil {
		q.tail = nil
	}
	return item.data
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
			return
		}
		prev = item
	}
}

func (q *Queue[T]) RemoveAll(data T) {
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
		} else {
			prev = item
		}
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
			return
		}
		prev = item
	}
}

func (q *Queue[T]) RemoveAllEx(data T, isEqual func(T, T) bool) {
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
		} else {
			prev = item
		}
	}
}
