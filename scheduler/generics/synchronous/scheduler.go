package synchronous

type LazyFunc[T permittedValues] func(n ...T) T

type Scheduler[T permittedValues] interface {
	Size() int
	Add(fn LazyFunc[T], args ...T)
	Run() []Result[T]
}

type permittedValues interface {
	~int | ~float64 | ~string
}

type Result[T permittedValues] struct {
	Value T
	Err   error
}

type lazyScheduler[T permittedValues] struct {
	jobs []job[T]
}

type job[T permittedValues] struct {
	lazyFn func(n ...T) T
	args   []T
}

func New[T permittedValues]() Scheduler[T] {
	return &lazyScheduler[T]{
		jobs: []job[T]{},
	}
}

func (s *lazyScheduler[T]) Size() int {
	return len(s.jobs)
}

func (s *lazyScheduler[T]) Add(fn LazyFunc[T], args ...T) {
	s.jobs = append(s.jobs, job[T]{lazyFn: fn, args: args})
}

func (s *lazyScheduler[T]) Run() []Result[T] {
	res := make([]Result[T], s.Size())

	for i, job := range s.jobs {
		r := job.lazyFn(job.args...)
		res[i] = Result[T]{Value: r}
	}
	return res
}
