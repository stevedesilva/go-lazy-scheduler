package synchronous

import (
	"strings"
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
	var funcMal LazyFunc[int] = func(num ...int) int {
		return 0
	}

	s := New[int]()
	s.Add(failFn, 1, 2, 3)
	s.Add(add, 1, 1, 1)
	s.Add(funcMal, 1, 1, 1)

	assert.Equal(t, s.Size(), 3)

	s.Add(funcMal, 2, 3, 4)

}

func TestLazyScheduler_ShouldExecuteFunctionsInScheduledOrderWithIntValues(t *testing.T) {
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

	s := New[int]()
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

func TestLazyScheduler_ShouldExecuteFunctionsInScheduledOrderWithFloatValues(t *testing.T) {
	add := func(n ...float64) float64 {
		var res float64
		for _, v := range n {
			res += v
		}
		return res
	}

	multiply := func(n ...float64) float64 {
		res := 1.0
		for _, v := range n {
			res *= v
		}
		return res
	}

	s := New[float64]()
	s.Add(add, 2, 4, 98)
	s.Add(add, 54, 22, 29)
	s.Add(multiply, 2, 2, 2, 2)
	s.Add(multiply, 3, 7)

	got := s.Run()

	assertion := assert.New(t)
	assertion.Equal(4, s.Size())
	assertion.Equal(104.0, got[0].Value)
	assertion.Equal(105.0, got[1].Value)
	assertion.Equal(16.0, got[2].Value)
	assertion.Equal(21.0, got[3].Value)

}

func TestLazyScheduler_ShouldExecuteFunctionsInScheduledOrderWithString(t *testing.T) {
	add := func(n ...string) string {
		return strings.Join(n, ",")
	}

	s := New[string]()
	s.Add(add, "One", "Two", "Three")
	s.Add(add, "do", "ray", "me")

	got := s.Run()

	assertion := assert.New(t)
	assertion.Equal(2, s.Size())
	assertion.Equal("One,Two,Three", got[0].Value)
	assertion.Equal("do,ray,me", got[1].Value)

}

//run `go test -bench=.` in
func BenchmarkLazyScheduler1000Jobs(b *testing.B) {
	LazySchedulerBenchmark(b, 1000)
}

func BenchmarkLazyScheduler10000Jobs(b *testing.B) {
	LazySchedulerBenchmark(b, 10000)
}

// capture results globally to avoid optimization
var results = make([]Result[int], 0)

func LazySchedulerBenchmark(b *testing.B, count int) {
	s := New[int]()
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
