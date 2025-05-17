package pool

type IPool[T any] interface {
	Put(obj T)
	Get() T
	Size() int
	Clear()
}

type Pool[T any] struct {
	base *basePool[T]
}

func NewPool[T any]() *Pool[T] {
	return &Pool[T]{
		base: newBasePool(DefaultConfig[T]()),
	}
}

func NewPoolWithConfig[T any](cfg Config[T]) *Pool[T] {
	return &Pool[T]{
		base: newBasePool(cfg),
	}
}

func (p *Pool[T]) Put(obj T) {
	p.base.put(obj)
}

func (p *Pool[T]) Get() T {
	return p.base.get()
}

func (p *Pool[T]) Size() int {
	return p.base.size()
}

func (p *Pool[T]) Clear() {
	p.base.clear()
}
