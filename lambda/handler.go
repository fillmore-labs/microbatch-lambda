package function

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/fillmore-labs/microbatch-lambda/api"
)

func Handler(
	_ context.Context,
	request *events.APIGatewayV2HTTPRequest,
) (*events.APIGatewayV2HTTPResponse, error) {
	jobs, err := UnmarshalBody(request)
	if err != nil {
		return nil, err
	}

	results := ProcessJobs(jobs)

	return MarshalResponse(results)
}

func UnmarshalBody(request *events.APIGatewayV2HTTPRequest) ([]*api.Job, error) {
	var err error
	var body []byte
	if request.IsBase64Encoded {
		body, err = base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			return nil, err
		}
	} else {
		body = []byte(request.Body)
	}

	var jobs []*api.Job
	err = json.Unmarshal(body, &jobs)

	return jobs, err
}

func MarshalResponse(results []*api.JobResult) (*events.APIGatewayV2HTTPResponse, error) {
	response, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	return &events.APIGatewayV2HTTPResponse{
		StatusCode:      http.StatusOK,
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            string(response),
		IsBase64Encoded: false,
	}, nil
}
