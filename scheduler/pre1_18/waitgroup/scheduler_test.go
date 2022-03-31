package synchronous_test

import (
	. "go-lazy-scheduler/scheduler/pre1_18/waitgroup"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScheduler_ShouldAddFunctionsLazily(t *testing.T) {
	add := func(n ...int) int {
		sum := 0
		for _, v := range n {
			sum += v
		}
		return sum
	}

	failFn := func(n ...int) int {
		t.Errorf("Fuction should not be called when added")
		return 0
	}

	s := New()
	s.Add(failFn, 1, 2, 3)
	s.Add(add, 1, 1, 1)

	assert.Equal(t, s.Size(), 2)
}

func TestLazyScheduler_ShouldExecuteFunctionsInScheduledOrder(t *testing.T) {
	add := func(n ...int) int {
		res := 0
		for _, v := range n {
			res += v
		}
		return res
	}

	multiply := func(n ...int) int {
		res := 1
		for _, v := range n {
			res *= v
		}
		return res
	}

	s := New()
	s.Add(add, 2, 4, 98)
	s.Add(add, 54, 22, 29)
	s.Add(multiply, 2, 2, 2, 2)
	s.Add(multiply, 3, 7)

	got := s.Run()

	assertion := assert.New(t)
	assertion.Equal(4, s.Size())
	assertion.Equal(104, got[0].Value)
	assertion.Equal(105, got[1].Value)
	assertion.Equal(16, got[2].Value)
	assertion.Equal(21, got[3].Value)

}

// run `go test -bench=.` in
func BenchmarkLazyScheduler1000Jobs(b *testing.B) {
	LazySchedulerBenchmark(b, 1000)
}

func BenchmarkLazyScheduler10000Jobs(b *testing.B) {
	LazySchedulerBenchmark(b, 10000)
}

// capture results globally to avoid optimization
var results = make([]Result, 0)

func LazySchedulerBenchmark(b *testing.B, count int) {
	s := New()
	add := func(n ...int) int {
		res := 0
		for _, v := range n {
			res += v
		}
		return res
	}
	for i := 0; i < count; i++ {
		s.Add(add, 5, 5, 5)
	}
	// Our terminal Go will pass a variable `N` into this function numerous times until it gets a relatively consistent result
	for i := 0; i < b.N; i++ {
		results = s.Run()
	}
}
