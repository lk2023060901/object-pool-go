package pool

type Config[T any] struct {
	New     func() T
	MaxSize int
	OnPut   func(obj T)
	OnGet   func(obj T)
}

const DefaultObjectMaxSize int = 64

func DefaultConfig[T any]() Config[T] {
	return Config[T]{
		New: func() T {
			var obj T
			return obj
		},
		MaxSize: DefaultObjectMaxSize,
	}
}
