package synchronous

type permittedValues interface {
	~int | ~float64 | ~string
}

//type Scheduler[T permittedValues] interface {
//	Size() int
//	Add(fn LazyFunc, args ...T)
//	Run() []Result[T]
//}

type Scheduler[T permittedValues] interface {
	Size() int
	Add(fn func(n ...T) T, args ...T)
	Run() []Result[T]
}

type Result[T permittedValues] struct {
	Value T
	Err   error
}

type LazyScheduler[T permittedValues] struct {
	jobs []job[T]
}

type LazyFunc[T permittedValues] func(n ...T) T

//type job[T permittedValues] struct {
//	lazyFn LazyFunc[T]
//	args   []T
//}

type job[T permittedValues] struct {
	lazyFn func(n ...T) T
	args   []T
}

func New[T permittedValues]() Scheduler[T] {
	return &LazyScheduler[T]{
		jobs: []job[T]{},
	}
}

//func New[T permittedValues]() *LazyScheduler[T] {
//	return &LazyScheduler[T]{
//		jobs: []job[T]{},
//	}
//}

func (s *LazyScheduler[T]) Size() int {
	return len(s.jobs)
}

//func (s *LazyScheduler[T]) Add(fn LazyFunc, args ...T) {
//	s.jobs = append(s.jobs, job[T]{lazyFn: fn, args: args})
//}

func (s *LazyScheduler[T]) Add(fn func(n ...T) T, args ...T) {
	s.jobs = append(s.jobs, job[T]{lazyFn: fn, args: args})
}

func (s *LazyScheduler[T]) Run() []Result[T] {
	res := make([]Result[T], s.Size())

	for i, job := range s.jobs {
		r := job.lazyFn(job.args...)
		res[i] = Result[T]{Value: r}
	}
	return res
}
