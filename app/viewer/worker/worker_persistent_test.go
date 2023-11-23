package worker

import (
	"darkbot/app/settings/logus"
	"darkbot/app/viewer/worker/worker_logus"
	"darkbot/app/viewer/worker/worker_types"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkerPersistent(t *testing.T) {
	jobPool := NewJobPool[*JobTest](
		WithAllowFailedJobs[*JobTest](),
		WithDisableParallelism[*JobTest](false),
	)

	jobs := []*JobTest{}
	for job_id := 1; job_id <= 3; job_id++ {
		jobs = append(jobs, NewJobTest(worker_types.JobID(job_id)))
	}

	jobPool.RunJobPool(jobs)

	done_count := 0
	failed_count := 0
	for job_number, job := range jobs {
		logus.Debug(fmt.Sprintf("job.Done=%t", job.done), worker_logus.JobNumber(worker_types.JobID(job_number)), JobResult(job.result))
		if job.done {
			done_count += 1
		} else {
			failed_count += 1
		}
	}
	assert.GreaterOrEqual(t, done_count, 3, "expected finding done jobs")
	assert.LessOrEqual(t, failed_count, 3, "expected finding failed jobs because of time sleep greater than timeout")
}
