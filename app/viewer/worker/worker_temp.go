package worker

import (
	"darkbot/app/settings/logus"
	"darkbot/app/viewer/worker/worker_logus"
	"darkbot/app/viewer/worker/worker_types"
	"time"
)

// ====================

type IJob interface {
	runJob(worker_id worker_types.WorkerID) worker_types.JobStatusCode
	isDone() bool
}

type Job struct {
	id   worker_types.JobID
	done bool
}

func NewJob(id worker_types.JobID) *Job {
	return &Job{id: id}
}

func (data *Job) isDone() bool { return data.done }

const (
	CodeSuccess worker_types.JobStatusCode = 0
	CodeTimeout worker_types.JobStatusCode = 1
)

type JobPool[jobd IJob] struct {
	jobTimeout worker_types.Seconds
	numWorkers int

	allow_failed_jobs   bool
	disable_parallelism worker_types.DebugDisableParallelism
}

type JobPoolOption[T IJob] func(r *JobPool[T])

func WithAllowFailedJobs[T IJob]() JobPoolOption[T] {
	return func(c *JobPool[T]) {
		c.allow_failed_jobs = true
	}
}

func WithWorkersAmount[T IJob](value int) JobPoolOption[T] {
	return func(c *JobPool[T]) { c.numWorkers = value }
}

func WithJobTimeout[T IJob](value worker_types.Seconds) JobPoolOption[T] {
	return func(c *JobPool[T]) { c.jobTimeout = value }
}

func WithDisableParallelism[T IJob](disable_parallelism worker_types.DebugDisableParallelism) JobPoolOption[T] {
	return func(c *JobPool[T]) { c.disable_parallelism = disable_parallelism }
}

func NewJobPool[T IJob](opts ...JobPoolOption[T]) *JobPool[T] {
	j := &JobPool[T]{
		numWorkers: 3,
		jobTimeout: 30,
	}

	for _, opt := range opts {
		opt(j)
	}

	return j
}

func (j *JobPool[jobd]) launchWorker(worker_id worker_types.WorkerID, jobs <-chan jobd, results chan<- worker_types.JobStatusCode) {
	logus.Debug("worker started", worker_logus.WorkerID(worker_id))
	for job := range jobs {
		results <- job.runJob(worker_id)
	}
	logus.Debug("worker finished", worker_logus.WorkerID(worker_id))
}

/// Temporal

func (j *JobPool[jobd]) runJobsinTemporalWorkers(jobs []jobd) []worker_types.JobStatusCode {
	numJobs := len(jobs)

	// In order to use our pool of workers we need to send them work and collect their results.
	// We make 2 channels for this.
	jobs_channel := make(chan jobd, numJobs)
	result_channel := make(chan worker_types.JobStatusCode, numJobs)
	status_codes := []worker_types.JobStatusCode{}

	// This starts up N workers, initially blocked because there are no jobs yet.
	for worker_id := 1; worker_id <= j.numWorkers; worker_id++ {
		go j.launchWorker(worker_types.WorkerID(worker_id), jobs_channel, result_channel)
	}

	// Here we send 5 jobs and
	for _, job := range jobs {
		jobs_channel <- job
	}

	// then close that channel to indicate that is all the work we have.
	close(jobs_channel)

	// Finally we collect all the results of the work.
	// This also ensures that the worker goroutines have finished.
	// An alternative way to wait for multiple goroutines is to use a WaitGroup.
	for job_number := range jobs {
		select {
		case res := <-result_channel:
			status_codes = append(status_codes, res)
		case <-time.After(time.Duration(j.jobTimeout) * time.Second):
			// non zero exit by timeout
			status_codes = append(status_codes, CodeTimeout)
			logus.Error("timeout for", worker_logus.JobNumber(worker_types.JobID(job_number)))
		}

	}
	return status_codes
}

func (jobPool *JobPool[jobd]) RunTemporalPool(jobs []jobd) {
	/*
		Switcher executing jobs with smth resembling multithreaded pool
		or executing synrconously if debug is on
	*/

	if jobPool.disable_parallelism {
		for pseudo_worker_id, job := range jobs {
			job.runJob(worker_types.WorkerID(pseudo_worker_id))
		}
	} else {
		status_codes := jobPool.runJobsinTemporalWorkers(jobs)
		logus.Debug("results", worker_logus.StatusCodes(status_codes))
	}

	for job_number, job := range jobs {
		job_id := worker_types.JobID(job_number)
		if !job.isDone() && !jobPool.allow_failed_jobs {
			logus.Error("job failed", worker_logus.JobNumber(job_id))
		}
		logus.Debug("job succeed", worker_logus.JobNumber(job_id))
	}
}
