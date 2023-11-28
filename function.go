package function

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/fillmore-labs/microbatch-lambda/pkg/api"
)

func ProcessJobs(jobs []api.Job) []api.JobResult {
	jobIDs := make([]string, 0, len(jobs))
	for _, job := range jobs {
		jobIDs = append(jobIDs, strconv.FormatInt(int64(job.ID), 10))
	}
	slog.Info("Processing", "jobs", strings.Join(jobIDs, ", "))

	results := make([]api.JobResult, 0, len(jobs))
	for _, job := range jobs {
		results = append(results, ProcessJob(job))
	}

	return results
}

func ProcessJob(job api.Job) api.JobResult {
	return api.JobResult{
		ID:   job.ID,
		Body: fmt.Sprintf("Hello, %s!", job.Body),
	}
}
