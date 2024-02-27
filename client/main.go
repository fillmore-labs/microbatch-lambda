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
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"fillmore-labs.com/microbatch"
	"github.com/aws/aws-sdk-go-v2/config"
	pb "github.com/fillmore-labs/microbatch-lambda/api/proto/v1alpha1"
	"github.com/spf13/viper"
)

const (
	iterations = 50
	delay      = 12 * time.Millisecond
	timeout    = 1500 * time.Millisecond
	batchSize  = 15
	batchDelay = 250 * time.Millisecond
)

var errExecutionTime = errors.New("execution time exceeded")

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	batcher, err := createBatcher(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create batcher", "error", err)

		return
	}

	slog.InfoContext(ctx, "Submitting jobs...")

	var wg sync.WaitGroup
	for i := 0; i < iterations; i++ {
		start := time.Now()
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			ctx2, cancel2 := context.WithTimeoutCause(ctx, timeout, errExecutionTime)
			defer cancel2()

			result, err := unwrap(batcher.Execute(ctx2, &pb.Job{
				Body:          fmt.Sprintf("Job %d", i),
				CorrelationId: int64(i),
			}))
			if err != nil {
				slog.WarnContext(ctx2, "Job failed to execute", "jobID", i, "error", err)

				return
			}

			slog.InfoContext(ctx2, "Job processed", "jobID", i, "result", result, "duration", time.Since(start))
		}(i)

		time.Sleep(delay)
	}

	batcher.Send()
	wg.Wait()

	slog.InfoContext(ctx, "Done...")
}

func createBatcher(ctx context.Context) (*microbatch.Batcher[*pb.Job, *pb.JobResult], error) {
	pcfg, err := readProcessorConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration: %w", err)
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	processor := NewRemoteProcessor(ctx, cfg.Credentials, pcfg)

	return microbatch.NewBatcher(
		processor.ProcessJobs,
		(*pb.Job).GetCorrelationId,
		(*pb.JobResult).GetCorrelationId,
		microbatch.WithSize(batchSize),
		microbatch.WithTimeout(batchDelay),
	), nil
}

func readProcessorConfig() (ProcessorConfig, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return ProcessorConfig{}, fmt.Errorf("can't read configuration file: %w", err)
	}

	var cfg ProcessorConfig
	err := v.UnmarshalExact(&cfg)

	return cfg, err
}
