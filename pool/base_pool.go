package pool

type basePool[T any] struct {
	cfg  Config[T]
	objs []T
}

func newBasePool[T any](cfg Config[T]) *basePool[T] {
	return &basePool[T]{
		cfg:  cfg,
		objs: make([]T, 0, cfg.MaxSize),
	}
}

func (p *basePool[T]) put(obj T) bool {
	if len(p.objs) >= p.cfg.MaxSize {
		return false
	}
	if p.cfg.OnPut != nil {
		p.cfg.OnPut(obj)
	}
	p.objs = append(p.objs, obj)
	return true
}

func (p *basePool[T]) get() T {
	n := len(p.objs)
	if n == 0 {
		if p.cfg.New != nil {
			return p.cfg.New()
		}
		var zero T
		return zero
	}
	obj := p.objs[n-1]
	p.objs = p.objs[:n-1]
	if p.cfg.OnGet != nil {
		p.cfg.OnGet(obj)
	}
	return obj
}

func (p *basePool[T]) size() int {
	return len(p.objs)
}

func (p *basePool[T]) clear() {
	p.objs = p.objs[:0]
}
