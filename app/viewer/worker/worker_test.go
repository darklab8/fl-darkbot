package worker

import (
	"darkbot/app/settings/logus"
	"darkbot/app/viewer/worker/worker_logus"
	"darkbot/app/viewer/worker/worker_types"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ======================
// Test Example

type JobTest struct {
	Job
	// any desired arbitary data
	result int
}

func (data *JobTest) runJob(worker_id worker_types.WorkerID) worker_types.JobStatusCode {
	// logus.Debug("", "worker", worker_id, "started  job", data.id)
	time.Sleep(time.Second * time.Duration(data.id))
	// logus.Debug("", "worker", worker_id, "finished job", data.id)
	data.result = data.id * 1
	data.done = true
	return CodeSuccess
}

func TestWorker(t *testing.T) {
	jobPool := NewJobPool[*JobTest](
		WithAllowFailedJobs[*JobTest](true),
	)

	jobs := []*JobTest{}
	for job_id := 1; job_id <= 3; job_id++ {
		jobs = append(jobs, &JobTest{Job: Job{id: job_id}})
	}

	RunJobPool(worker_types.DebugDisableParallelism(false), jobPool, jobs)

	done_count := 0
	failed_count := 0
	for job_number, job := range jobs {
		logus.Debug(fmt.Sprintf("job.Done=%t", job.done), worker_logus.JobNumber(job_number), worker_logus.JobResult(job.result))
		if job.done {
			done_count += 1
		} else {
			failed_count += 1
		}
	}
	assert.GreaterOrEqual(t, done_count, 3, "expected finding done jobs")
	assert.LessOrEqual(t, failed_count, 3, "expected finding failed jobs because of time sleep greater than timeout")
}
