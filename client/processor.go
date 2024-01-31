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
)

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

	return p.parseResponse(response)
}

type ProcessorConfig struct {
	LambdaURL string        `mapstructure:"lambda-url"`
	Region    string        `mapstructure:"region"`
	Timeout   time.Duration `mapstructure:"timeout"`
}

func NewRemoteProcessor(
	ctx context.Context,
	credentials aws.CredentialsProvider,
	cfg ProcessorConfig,
) *RemoteProcessor {
	return &RemoteProcessor{
		// Pass the context to the processor, https://go.dev/blog/context-and-structs
		ctx:         ctx,
		credentials: credentials,
		httpSigner:  signer.NewSigner(signerOptions),
		region:      cfg.Region,
		lambdaURL:   cfg.LambdaURL,
		timeOut:     cfg.Timeout,
	}
}

func signerOptions(s *signer.SignerOptions) {
	s.DisableURIPathEscaping = true
}
