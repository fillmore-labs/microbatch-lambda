package function

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	pb "github.com/fillmore-labs/microbatch-lambda/api/proto/v1alpha1"
	"google.golang.org/protobuf/proto"
)

func Handler(
	ctx context.Context,
	request *events.APIGatewayV2HTTPRequest,
) (*events.APIGatewayV2HTTPResponse, error) {
	jobs, err := UnmarshalBody(request)
	if err != nil {
		return nil, err
	}

	results, err := ProcessJobs(ctx, jobs)
	if err != nil {
		return nil, err
	}

	return MarshalResponse(results)
}

func UnmarshalBody(request *events.APIGatewayV2HTTPRequest) (Jobs, error) {
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

	var requests pb.BatchRequest
	err = proto.Unmarshal(body, &requests)

	return requests.GetJobs(), err
}

func MarshalResponse(results JobResults) (*events.APIGatewayV2HTTPResponse, error) {
	msg := &pb.BatchResponse{Results: results}

	response, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	messageType := msg.ProtoReflect().Descriptor().FullName()
	contentType := fmt.Sprintf("application/x-protobuf; messageType=\"%s\"", messageType)

	return &events.APIGatewayV2HTTPResponse{
		StatusCode:      http.StatusOK,
		Headers:         map[string]string{"Content-Type": contentType},
		Body:            string(response),
		IsBase64Encoded: false,
	}, nil
}
