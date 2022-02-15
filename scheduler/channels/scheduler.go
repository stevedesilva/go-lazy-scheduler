package scheduler

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


type jobResult struct {
	idx   int
	value int
	err   error
}
type resultChan chan jobResult

func (s *LazyScheduler) Run() []Result {
	res := make([]Result, s.Size())
	out := make(resultChan, s.Size())
	defer func() {
		close(out)
	}()

	for i, j := range s.jobs {
		go func(idx int, job job) {
			val := job.lazyFn(job.args...)
			res := jobResult{idx: idx, value: val, err: nil}
			out <- res
		}(i, j)
		// when launching goroutines in a range loop we need to copy the values otherwise the closure will simply
		// use the last i, j pair from the range loop
	}

	for i := 0; i < s.Size(); i++ {
		r := <- out
		res[r.idx] = Result{Value: r.value, Err: r.err}
	}

	return res
}
