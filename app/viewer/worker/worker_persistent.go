package worker

import (
	"darkbot/app/settings/logus"
	"darkbot/app/viewer/worker/worker_logus"
	"darkbot/app/viewer/worker/worker_types"
	"time"
)

type JobPoolPeristent[jobd IJob] struct {
	*JobPool[jobd]

	jobs_channel   chan jobd
	result_channel chan worker_types.JobStatusCode
}

func NewJobPoolPersistent[jobd IJob](opts ...JobPoolOption[jobd]) *JobPoolPeristent[jobd] {
	j := &JobPoolPeristent[jobd]{JobPool: NewJobPool[jobd](opts...)}

	j.jobs_channel = make(chan jobd)
	j.result_channel = make(chan worker_types.JobStatusCode)

	// This starts up N workers, initially blocked because there are no jobs yet.
	for worker_id := 1; worker_id <= j.numWorkers; worker_id++ {
		go j.launchWorker(
			worker_types.WorkerID(worker_id),
			j.jobs_channel,
			j.result_channel,
		)
	}
	return j
}

func (j *JobPoolPeristent[jobd]) DelayJob(job jobd) {
	j.jobs_channel <- job
}

func (j *JobPoolPeristent[jobd]) AwaitSomeJob() {
	select {
	case status_code := <-j.result_channel:
		logus.Debug("finished some job succesfully", worker_logus.StatusCode(status_code))
	case <-time.After(time.Duration(j.jobTimeout) * time.Second):
		// non zero exit by timeout
		logus.Error("finished jobs with", worker_logus.StatusCode(CodeTimeout)) // TODO add worker_logus.JobNumber(worker_types.JobID(job_number)
	}
}

func (j *JobPoolPeristent[jobd]) AwaitTimeouts() {
	for {
		j.AwaitSomeJob()
	}
}
