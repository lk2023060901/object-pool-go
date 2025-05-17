package pool

import "sync"

type ConcurrencyPool[T any] struct {
	base *basePool[T]
	mu   sync.RWMutex
}

func NewConcurrencyPool[T any]() *ConcurrencyPool[T] {
	return &ConcurrencyPool[T]{
		base: newBasePool(DefaultConfig[T]()),
	}
}

func NewConcurrencyPoolWithConfig[T any](cfg Config[T]) *ConcurrencyPool[T] {
	return &ConcurrencyPool[T]{
		base: newBasePool(cfg),
	}
}

func (p *ConcurrencyPool[T]) Put(obj T) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.base.put(obj)
}

func (p *ConcurrencyPool[T]) Get() T {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.base.get()
}

func (p *ConcurrencyPool[T]) Size() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.base.size()
}

func (p *ConcurrencyPool[T]) Clear() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.base.clear()
}
