package function

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	pb "github.com/fillmore-labs/microbatch-lambda/api/proto/v1alpha1"
)

type (
	Jobs       []*pb.Job
	JobResults []*pb.JobResult
)

// ProcessJobs processes a batch of jobs.
//
// This function is called by the lambda function.
//
// The function returns a list of results, each result corresponds to a job in the input list.
// If the job was processed successfully, the result will contain the result of the job.
// If the job failed, the result will contain an error message.
func ProcessJobs(ctx context.Context, jobs Jobs) (JobResults, error) {
	slog.Info("Processing", "jobs", JobIDs(jobs))

	results := make(JobResults, 0, len(jobs))
	for _, job := range jobs {
		result := &pb.JobResult{CorrelationId: job.GetCorrelationId()}

		if body, err := ProcessJob(ctx, job.GetBody()); err == nil {
			result.Result = &pb.JobResult_Body{Body: body}
		} else {
			result.Result = &pb.JobResult_Error{Error: err.Error()}
		}

		results = append(results, result)
	}

	return results, nil
}

// ProcessJob processes a single job.
func ProcessJob(_ context.Context, body string) (string, error) {
	return fmt.Sprintf("Hello, %s!", body), nil
}

type JobIDs []*pb.Job

// LogValue implements [slog.LogValuer].
func (jobs JobIDs) LogValue() slog.Value {
	var b []byte
	if len(jobs) > 0 {
		b = strconv.AppendInt(b, jobs[0].GetCorrelationId(), 10) //nolint:gomnd
		for _, j := range jobs[1:] {
			b = append(b, ", "...)
			b = strconv.AppendInt(b, j.GetCorrelationId(), 10) //nolint:gomnd
		}
	}

	return slog.StringValue(string(b))
}
