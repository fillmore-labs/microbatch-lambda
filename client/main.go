package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"fillmore-labs.com/exp/async"
	"fillmore-labs.com/microbatch"
	"github.com/aws/aws-sdk-go-v2/config"
	pb "github.com/fillmore-labs/microbatch-lambda/api/proto/v1alpha1"
	"github.com/spf13/viper"
)

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

func main() {
	pcfg, err := readProcessorConfig()
	if err != nil {
		log.Fatalf("Failed to read configuration: %v", err)
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	processor := NewRemoteProcessor(ctx, cfg.Credentials, pcfg)

	const batchSize = 15
	const batchDelay = 250 * time.Millisecond

	batcher := microbatch.NewBatcher(
		processor.ProcessJobs,
		(*pb.Job).GetCorrelationId,
		(*pb.JobResult).GetCorrelationId,
		microbatch.WithSize(batchSize),
		microbatch.WithTimeout(batchDelay),
	)

	log.Println("Submitting jobs...")

	const iterations = 50
	const delay = 12 * time.Millisecond

	var wg sync.WaitGroup
	for i := 0; i < iterations; i++ {
		request := &pb.Job{
			Body:          fmt.Sprintf("Job %d", i),
			CorrelationId: int64(i),
		}

		future := batcher.SubmitJob(request)

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if result, err := async.Then(ctx, future, unwrap); err == nil {
				log.Printf("Result of job %d: %s\n", i, result)
			} else {
				log.Printf("Error executing job %d: %v\n", i, err)
			}
		}(i)

		time.Sleep(delay)
	}
	batcher.Shutdown()

	wg.Wait()

	log.Println("Done...")
}

type remoteError struct {
	msg string
}

func (r *remoteError) Error() string {
	return r.msg
}

var errMissingResult = &remoteError{"missing result"}

func unwrap(result *pb.JobResult) (string, error) {
	r := result.GetResult()

	switch r := r.(type) {
	case *pb.JobResult_Body:
		return r.Body, nil

	case *pb.JobResult_Error:
		return "", &remoteError{r.Error}

	default:
		return "", errMissingResult
	}
}
