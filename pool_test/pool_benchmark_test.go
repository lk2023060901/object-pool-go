package pool_test

import (
	"github.com/lk2023060901/object-pool-go/pool"
	"testing"
)

type LargeObject struct {
	data [1024]byte
}

// 防止优化用的变量
var sink *LargeObject

// 不使用对象池，直接新建对象（真实分配）
func Benchmark_NewWithoutPool(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		obj := &LargeObject{}
		sink = obj
	}
}

// 非并发对象池复用
func Benchmark_PoolReuse(b *testing.B) {
	p := pool.NewPoolWithConfig[*LargeObject](pool.Config[*LargeObject]{
		New:     func() *LargeObject { return &LargeObject{} },
		MaxSize: 1024,
	})
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		obj := p.Get()
		sink = obj
		p.Put(obj)
	}
}

// 并发安全对象池复用
func Benchmark_ConcurrencyPoolReuse(b *testing.B) {
	p := pool.NewConcurrencyPoolWithConfig[*LargeObject](pool.Config[*LargeObject]{
		New:     func() *LargeObject { return &LargeObject{} },
		MaxSize: 1024,
	})
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			obj := p.Get()
			sink = obj
			p.Put(obj)
		}
	})
}

// 并发场景下直接新建对象（真实分配）
func Benchmark_NewWithoutPool_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			obj := &LargeObject{}
			sink = obj
		}
	})
}
