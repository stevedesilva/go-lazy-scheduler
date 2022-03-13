package scheduler

import (
	"fmt"
	"github.com/google/uuid"
)

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

func NewWithGr(gr int) Scheduler {
	return &LazyScheduler{
		gr: gr,
	}
}

type LazyScheduler struct {
	jobs []job
	gr   int
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
	quit := make(chan struct{})
	defer close(quit)
	defer close(out)

	// if gr set then limit
	if s.gr > 0 {
		workers := make([]workerJobCh, 0)

		for i, j := range s.jobs {
			// create n workers
			if len(workers) < s.gr {
				fmt.Printf("i-%d creating ...\n", i)

				workers = append(workers, createWorker(out, quit))
			}

			// assign worker a job round robin
			count := i % len(workers)
			fmt.Printf("i-%d Count-%d Workers# %d \n", i, count, len(workers))
			workers[count] <- jobData{
				idx: i,
				job: j,
			}

		}

	} else {
		for i, j := range s.jobs {
			// no limit
			go func(idx int, job job) {
				val := job.lazyFn(job.args...)
				res := jobResult{idx: idx, value: val, err: nil}
				out <- res
			}(i, j)
			// when launching goroutines in a range loop we need to copy the values otherwise the closure will simply
			// use the last i, j pair from the range loop
		}
	}

	for i := 0; i < s.Size(); i++ {
		r := <-out
		res[r.idx] = Result{Value: r.value, Err: r.err}
	}

	return res
}

type jobData struct {
	idx int
	job
}

// read only chang
type workerJobCh chan jobData

func createWorker(out resultChan, exit chan struct{}) workerJobCh {
	in := make(workerJobCh)
	workerId := uuid.New().String()
	fmt.Printf("Wrk number %s: created \n", workerId)
	go func(wid string) {
		for {
			select {
			case data := <-in:
				fmt.Printf("Wrk number %s: processing job %v \n", workerId, data)
				val := data.lazyFn(data.args...)
				res := jobResult{idx: data.idx, value: val, err: nil}
				out <- res
			case <-exit:
				return
			}
		}
	}(workerId)
	return in
}
