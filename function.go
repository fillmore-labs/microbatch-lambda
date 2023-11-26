package function

import (
	"encoding/json"
	"fmt"

	"github.com/fillmore-labs/microbatch-lambda/pkg/api"
)

func ProcessJobs(s []byte) ([]byte, error) {
	var j []api.Job
	err := json.Unmarshal(s, &j)
	if err != nil {
		return nil, err
	}

	r := make([]api.JobResult, 0, len(j))
	for _, job := range j {
		id := job.CorrelationID()
		r = append(r, api.JobResult{
			ID:   id,
			Body: fmt.Sprintf("Hello, %s!", job.Body),
		})
	}

	rr, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return rr, nil
}
