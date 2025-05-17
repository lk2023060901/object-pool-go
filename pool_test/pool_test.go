package pool_test

import (
	"github.com/lk2023060901/object-pool-go/pool"
	"sync/atomic"
	"testing"
)

func TestPoolBasic(t *testing.T) {
	p := pool.NewPoolWithConfig[int](pool.Config[int]{
		New:     func() int { return 999 },
		MaxSize: 2,
	})

	if val := p.Get(); val != 999 {
		t.Errorf("expected default new value 999, got %d", val)
	}

	p.Put(1)
	p.Put(2)
	p.Put(3) // 超过 MaxSize，不应存入

	if p.Size() != 2 {
		t.Errorf("expected pool size 2, got %d", p.Size())
	}

	a := p.Get()
	b := p.Get()
	c := p.Get()

	if a != 2 || b != 1 || c != 999 {
		t.Errorf("expected LIFO 2,1,999 got %d,%d,%d", a, b, c)
	}

	if p.Size() != 0 {
		t.Errorf("expected pool size 0 after Get, got %d", p.Size())
	}
}

func TestPool_OnPut_OnGet(t *testing.T) {
	var putCount int32
	var getCount int32
	var log []string

	p := pool.NewPoolWithConfig[string](pool.Config[string]{
		New:     func() string { return "new" },
		MaxSize: 5,
		OnPut: func(obj string) {
			atomic.AddInt32(&putCount, 1)
			log = append(log, "put:"+obj)
		},
		OnGet: func(obj string) {
			atomic.AddInt32(&getCount, 1)
			log = append(log, "get:"+obj)
		},
	})

	p.Put("a")
	p.Put("b")

	_ = p.Get()
	_ = p.Get()
	_ = p.Get() // will call New()

	if putCount != 2 || getCount != 2 {
		t.Errorf("expected putCount=2, getCount=2, got %d %d", putCount, getCount)
	}

	expectedLog := []string{"put:a", "put:b", "get:b", "get:a"}
	if len(log) != len(expectedLog) {
		t.Fatalf("log length mismatch: got %v", log)
	}
	for i, v := range expectedLog {
		if log[i] != v {
			t.Errorf("log mismatch at %d: expected %s, got %s", i, v, log[i])
		}
	}
}

func TestPool_Clear(t *testing.T) {
	p := pool.NewPool[int]()
	p.Put(1)
	p.Put(2)

	if p.Size() != 2 {
		t.Fatalf("expected size 2, got %d", p.Size())
	}

	p.Clear()

	if p.Size() != 0 {
		t.Errorf("expected cleared size 0, got %d", p.Size())
	}

	val := p.Get()
	if val != 0 {
		t.Errorf("expected zero value after clear, got %d", val)
	}
}
