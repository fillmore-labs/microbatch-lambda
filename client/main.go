package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

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
	ctx := context.Background()

	pcfg, err := readProcessorConfig()
	if err != nil {
		log.Fatalf("Failed to read configuration: %v", err)
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
	}

	requestContext, cancel := context.WithCancel(ctx)

	processor := NewRemoteProcessor(requestContext, cfg.Credentials, pcfg)

	const batchSize = 15
	const batchDelay = 250 * time.Millisecond

	//

	batcher := microbatch.NewBatcher(
		processor,
		(*pb.Job).GetCorrelationId,
		(*pb.JobResult).GetCorrelationId,
		microbatch.WithSize(batchSize),
		microbatch.WithTimeout(batchDelay),
	)

	log.Println("Submitting jobs...")

	const iterations = 50
	const delay = 12 * time.Millisecond

	var wg sync.WaitGroup
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		time.Sleep(delay)
		go submitWork(requestContext, &wg, batcher, int64(i+1))
	}
	wg.Wait()
	cancel()

	batcher.Shutdown()

	log.Println("Done...")
}

func submitWork(ctx context.Context, wg *sync.WaitGroup, batcher *microbatch.Batcher[*pb.Job, *pb.JobResult], i int64) {
	defer wg.Done()

	request := &pb.Job{
		Body:          fmt.Sprintf("Job %d", i),
		CorrelationId: i,
	}

	reply, err := batcher.ExecuteJob(ctx, request)

	if err != nil {
		log.Printf("Error executing job %d: %v\n", i, err)
	} else {
		result := reply.GetResult()

		switch r := result.(type) {
		case *pb.JobResult_Error:
			log.Printf("Error executing job %d: %s\n", i, r.Error)

		case *pb.JobResult_Body:
			log.Printf("Result of job %d: %s\n", i, r.Body)

		default:
			log.Printf("Missing result for job %d\n", i)
		}
	}
}
