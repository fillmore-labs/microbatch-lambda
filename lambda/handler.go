// Copyright 2023-2024 Oliver Eikemeier. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

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
