package pool_test

import (
	"github.com/lk2023060901/object-pool-go/pool"
	"sync"
	"sync/atomic"
	"testing"
)

func TestConcurrencyPoolBasic(t *testing.T) {
	p := pool.NewConcurrencyPoolWithConfig[int](pool.Config[int]{
		New:     func() int { return 100 },
		MaxSize: 2,
	})

	p.Put(1)
	p.Put(2)
	p.Put(3) // 超出 MaxSize

	if p.Size() != 2 {
		t.Errorf("expected size 2, got %d", p.Size())
	}

	a := p.Get()
	b := p.Get()
	c := p.Get()

	if a != 2 || b != 1 || c != 100 {
		t.Errorf("expected LIFO 2,1,100, got %d,%d,%d", a, b, c)
	}

	if p.Size() != 0 {
		t.Errorf("expected size 0, got %d", p.Size())
	}
}

func TestConcurrencyPool_OnPut_OnGet(t *testing.T) {
	var putCount int32
	var getCount int32

	p := pool.NewConcurrencyPoolWithConfig[string](pool.Config[string]{
		New:     func() string { return "new" },
		MaxSize: 3,
		OnPut: func(obj string) {
			atomic.AddInt32(&putCount, 1)
		},
		OnGet: func(obj string) {
			atomic.AddInt32(&getCount, 1)
		},
	})

	p.Put("x")
	p.Put("y")
	_ = p.Get()
	_ = p.Get()
	_ = p.Get() // 会调用 New()

	if putCount != 2 || getCount != 2 {
		t.Errorf("expected putCount=2, getCount=2, got %d, %d", putCount, getCount)
	}
}

func TestConcurrencyPool_ParallelAccess(t *testing.T) {
	p := pool.NewConcurrencyPoolWithConfig[int](pool.Config[int]{
		New:     func() int { return 0 },
		MaxSize: 1000,
	})

	var wg sync.WaitGroup
	const numWorkers = 100
	const opsPerWorker = 100

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for j := 0; j < opsPerWorker; j++ {
				p.Put(base + j)
				_ = p.Get()
			}
		}(i * 1000)
	}

	wg.Wait()

	// 最多应该不超过 MaxSize
	if size := p.Size(); size > 1000 {
		t.Errorf("expected size <= 1000, got %d", size)
	}
}

func TestConcurrencyPool_Clear(t *testing.T) {
	p := pool.NewConcurrencyPool[int]()
	p.Put(10)
	p.Put(20)

	if p.Size() != 2 {
		t.Fatalf("expected size 2, got %d", p.Size())
	}

	p.Clear()

	if p.Size() != 0 {
		t.Errorf("expected size 0 after clear, got %d", p.Size())
	}

	val := p.Get()
	if val != 0 {
		t.Errorf("expected default zero value after clear, got %d", val)
	}
}
