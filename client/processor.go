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

	"github.com/aws/aws-sdk-go-v2/aws"
	signer "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	pb "github.com/fillmore-labs/microbatch-lambda/api/proto/v1alpha1"
	"google.golang.org/protobuf/proto"
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

var ErrNotOk = errors.New("non-200 status")

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

const service = "lambda"

func (p *RemoteProcessor) createRequest(ctx context.Context, jobs Jobs) (*http.Request, error) {
	msg := &pb.BatchRequest{Jobs: jobs}

	body, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	request, err := p.newRequest(ctx, body)
	if err != nil {
		return nil, err
	}

	messageType := msg.ProtoReflect().Descriptor().FullName()
	contentType := fmt.Sprintf("application/x-protobuf; messageType=\"%s\"", messageType)
	request.Header.Set("Content-Type", contentType)

	err = p.signRequest(ctx, body, request)

	return request, err
}

func (p *RemoteProcessor) newRequest(ctx context.Context, body []byte) (*http.Request, error) {
	bodyReader := bytes.NewReader(body)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, p.lambdaURL, bodyReader)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func (p *RemoteProcessor) signRequest(ctx context.Context, body []byte, request *http.Request) error {
	credentials, err := p.credentials.Retrieve(ctx)
	if err != nil {
		return err
	}

	hash := sha256.Sum256(body)
	payloadHash := hex.EncodeToString(hash[:])

	return p.httpSigner.SignHTTP(ctx, credentials, request, payloadHash, service, p.region, time.Now())
}

func (p *RemoteProcessor) parseResponse(response *http.Response) (JobResults, error) {
	body, err := io.ReadAll(response.Body)
	_ = response.Body.Close()
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		slog.Warn("Response Error", "status", response.Status)

		return nil, fmt.Errorf("response code %d: %w", response.StatusCode, ErrNotOk)
	}

	var result pb.BatchResponse
	err = proto.Unmarshal(body, &result)

	return result.GetResults(), err
}
