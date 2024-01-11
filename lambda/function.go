package function

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	pb "github.com/fillmore-labs/microbatch-lambda/api/proto/v1alpha1"
)

type (
	Jobs       []*pb.Job
	JobResults []*pb.JobResult
)

func ProcessJobs(jobs Jobs) JobResults {
	slog.Info("Processing", "jobs", JobIDs(jobs))

	results := make(JobResults, 0, len(jobs))
	for _, job := range jobs {
		results = append(results, ProcessJob(job))
	}

	return results
}

func ProcessJob(job *pb.Job) *pb.JobResult {
	return &pb.JobResult{
		Result: &pb.JobResult_Body{
			Body: fmt.Sprintf("Hello, %s!", job.GetBody()),
		},
		CorrelationId: job.GetCorrelationId(),
	}
}

type JobIDs []*pb.Job

func (jobs JobIDs) LogValue() slog.Value {
	var b strings.Builder
	if len(jobs) > 0 {
		b.WriteString(strconv.FormatInt(jobs[0].GetCorrelationId(), 10))
		for _, j := range jobs[1:] {
			b.WriteString(", ")
			b.WriteString(strconv.FormatInt(j.GetCorrelationId(), 10))
		}
	}

	return slog.StringValue(b.String())
}
