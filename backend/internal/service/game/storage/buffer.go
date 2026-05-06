package storage

import (
	"errors"
)

type RingBuffer[T any] struct {
	tail int
	head int
	data []T
	cap  int
}

var (
	FullBuffer  = errors.New("ring buffer is full")
	EmptyBuffer = errors.New("ring buffer is empty")
)

func NewRingBuffer[T any](cap int) *RingBuffer[T] {
	return &RingBuffer[T]{
		data: make([]T, cap),
		cap:  cap,
	}
}

func (r *RingBuffer[T]) RemainingItemsNumber() int {
	return (r.tail - r.head + r.cap) % r.cap
}

func (r *RingBuffer[T]) Push(data T) error {
	if (r.tail+1)%r.cap == r.head {
		return FullBuffer
	}
	r.data[r.tail] = data
	r.tail = (r.tail + 1) % r.cap
	return nil
}

func (r *RingBuffer[T]) Pop() (T, error) {
	var zero T
	if r.tail == r.head {
		return zero, EmptyBuffer
	}
	res := r.data[r.head]
	r.head = (r.head + 1) % r.cap
	return res, nil
}
