package function

import (
	"context"
	"encoding/base64"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// Handler is the main handler for the function.
func Handler(
	_ context.Context,
	request *events.APIGatewayV2HTTPRequest,
) (*events.APIGatewayV2HTTPResponse, error) {
	var body []byte
	if request.IsBase64Encoded {
		b, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			return nil, err
		}
		body = b
	} else {
		body = []byte(request.Body)
	}

	response, err := ProcessJobs(body)
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
