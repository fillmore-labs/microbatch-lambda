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

package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	pb "github.com/fillmore-labs/microbatch-lambda/api/proto/v1alpha1"
	"google.golang.org/protobuf/proto"
)

// createRequest creates a HTTP request to AWS Lambda with an AWS signature version 4.
func (p *RemoteProcessor) createRequest(ctx context.Context, jobs Jobs) (*http.Request, error) {
	msg := &pb.BatchRequest{Jobs: jobs}

	body, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, p.lambdaURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", contentType(msg))

	err = p.signRequest(ctx, hash(body), request)

	return request, err
}

// contentType returns the content type of the protobuf message.
func contentType(msg proto.Message) string {
	messageType := msg.ProtoReflect().Descriptor().FullName()

	return fmt.Sprintf("application/x-protobuf; messageType=\"%s\"", messageType)
}

// hash calculates the SHA 256 hash of the body.
func hash(body []byte) string {
	hash := sha256.Sum256(body)

	return hex.EncodeToString(hash[:])
}

const service = "lambda"

// signRequest signs the request with AWS SigV4.
func (p *RemoteProcessor) signRequest(ctx context.Context, payloadHash string, request *http.Request) error {
	credentials, err := p.credentials.Retrieve(ctx)
	if err != nil {
		return err
	}

	return p.httpSigner.SignHTTP(ctx, credentials, request, payloadHash, service, p.region, time.Now())
}

var ErrNotOk = errors.New("non-200 status")

// parseResponse parses the HTTP response from AWS Lambda.
func (p *RemoteProcessor) parseResponse(ctx context.Context, response *http.Response) (JobResults, error) {
	if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()

		slog.WarnContext(ctx, "Response Error", "status", response.Status)

		return nil, fmt.Errorf("response code %d: %w", response.StatusCode, ErrNotOk)
	}

	body, err := io.ReadAll(response.Body)
	_ = response.Body.Close()
	if err != nil {
		return nil, err
	}

	var result pb.BatchResponse
	err = proto.Unmarshal(body, &result)

	return result.GetResults(), err
}
