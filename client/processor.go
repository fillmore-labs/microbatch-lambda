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
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	signer "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	pb "github.com/fillmore-labs/microbatch-lambda/api/proto/v1alpha1"
)

type (
	Jobs       []*pb.Job
	JobResults []*pb.JobResult

	RemoteProcessor struct {
		ctx         context.Context //nolint:containedctx
		credentials aws.CredentialsProvider
		httpSigner  signer.HTTPSigner
		region      string
		lambdaURL   string
		timeOut     time.Duration
	}

	ProcessorConfig struct {
		LambdaURL string        `mapstructure:"lambda-url"`
		Region    string        `mapstructure:"region"`
		Timeout   time.Duration `mapstructure:"timeout"`
	}
)

func NewRemoteProcessor(
	ctx context.Context,
	credentials aws.CredentialsProvider,
	cfg ProcessorConfig,
) *RemoteProcessor {
	signerOptions := func(s *signer.SignerOptions) {
		s.DisableURIPathEscaping = true
	}
	httpSigner := signer.NewSigner(signerOptions)

	return &RemoteProcessor{
		// Pass the context to the processor, https://go.dev/blog/context-and-structs
		ctx:         ctx,
		credentials: credentials,
		httpSigner:  httpSigner,
		region:      cfg.Region,
		lambdaURL:   cfg.LambdaURL,
		timeOut:     cfg.Timeout,
	}
}

// ProcessJobs sends the jobs to the remote lambda and returns the results.
func (p *RemoteProcessor) ProcessJobs(jobs Jobs) (JobResults, error) {
	ctx, cancel := context.WithTimeout(p.ctx, p.timeOut)
	defer cancel()

	request, err := p.createRequest(ctx, jobs)
	if err != nil {
		return nil, fmt.Errorf("can't build request %w", err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("can't send request %w", err)
	}

	return p.parseResponse(ctx, response)
}
