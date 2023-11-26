package api

type (
	JobID int64
	Job   struct {
		ID   JobID  `json:"id"`
		Body string `json:"body"`
	}
	JobResult struct {
		ID   JobID  `json:"id"`
		Body string `json:"body,omitempty"`
		Err  string `json:"error,omitempty"`
	}
)

func (j *Job) CorrelationID() JobID {
	return j.ID
}

func (j *JobResult) CorrelationID() JobID {
	return j.ID
}
