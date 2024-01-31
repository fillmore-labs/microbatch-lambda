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
func (p *RemoteProcessor) parseResponse(response *http.Response) (JobResults, error) {
	if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()

		slog.Warn("Response Error", "status", response.Status)

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
