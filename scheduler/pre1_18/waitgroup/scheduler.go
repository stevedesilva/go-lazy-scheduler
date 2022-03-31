package synchronous

import "sync"

type Scheduler interface {
	Size() int
	Add(fn lazyFunc, args ...int)
	Run() []Result
}

type Result struct {
	Value int
	Err   error
}

type lazyFunc func(n ...int) int

type job struct {
	lazyFn lazyFunc
	args   []int
}

func New() Scheduler {
	return &LazyScheduler{}
}

type LazyScheduler struct {
	jobs []job
}

func (s *LazyScheduler) Size() int {
	return len(s.jobs)
}

func (s *LazyScheduler) Add(fn lazyFunc, args ...int) {
	s.jobs = append(s.jobs, job{lazyFn: fn, args: args})
}

func (s *LazyScheduler) Run() []Result {
	wg := sync.WaitGroup{}
	wg.Add(s.Size())

	res := make([]Result, s.Size())

	for i, j := range s.jobs {
		go func(index int, job job) {
			r := job.lazyFn(job.args...)
			res[index] = Result{Value: r}
			defer wg.Done()
		}(i, j)
		// when launching goroutines in a range loop we need to copy the values otherwise the closure will simply
		// use the last i, j pair from the range loop
	}

	wg.Wait()
	return res
}
