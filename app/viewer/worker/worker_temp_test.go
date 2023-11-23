package worker

import (
	"darkbot/app/settings/logus"
	"darkbot/app/viewer/worker/worker_logus"
	"darkbot/app/viewer/worker/worker_types"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ======================
// Test Example

type JobTest struct {
	*Job

	// any desired arbitary data
	result worker_types.JobID
}

func NewJobTest(id worker_types.JobID) *JobTest {
	return &JobTest{Job: NewJob(id)}
}

func (data *JobTest) runJob(worker_id worker_types.WorkerID) worker_types.JobStatusCode {
	logus.Debug("job test started", worker_logus.WorkerID(worker_id), worker_logus.JobNumber(data.id))
	time.Sleep(time.Second * time.Duration(data.id))
	data.result = data.id * 1
	data.done = true
	logus.Debug("job test finished", worker_logus.WorkerID(worker_id), worker_logus.JobNumber(data.id))
	return CodeSuccess
}

func JobResult(value worker_types.JobID) logus.SlogParam {
	return func(c *logus.SlogGroup) {
		c.Params["job_result"] = strconv.Itoa(int(value))
	}
}

func TestWorkerTemp(t *testing.T) {
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
