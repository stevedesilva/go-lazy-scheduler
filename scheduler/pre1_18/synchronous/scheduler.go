package synchronous

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
	res := make([]Result, s.Size())

	for i, job := range s.jobs {
		r := job.lazyFn(job.args...)
		res[i] = Result{Value: r}
	}
	return res
}
