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
