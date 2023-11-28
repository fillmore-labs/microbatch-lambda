package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"fillmore-labs.com/microbatch"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/fillmore-labs/microbatch-lambda/pkg/api"
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

	var config ProcessorConfig
	err := v.UnmarshalExact(&config)

	return config, err
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

	batcher := microbatch.NewBatcher(
		processor,
		func(j *api.Job) api.JobID { return j.ID },
		func(r *api.JobResult) api.JobID { return r.ID },
		batchSize,
		batchDelay,
	)

	log.Println("Submitting jobs...")

	const iterations = 50
	const delay = 12 * time.Millisecond

	var wg sync.WaitGroup
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		time.Sleep(delay)
		go submitWork(requestContext, i, batcher, &wg)
	}
	wg.Wait()
	cancel()
	batcher.Shutdown()
	log.Println("Done...")
}

func submitWork(ctx context.Context, i int, batcher *microbatch.Batcher[*api.Job, *api.JobResult], wg *sync.WaitGroup) {
	request := &api.Job{
		ID:   api.JobID(i),
		Body: fmt.Sprintf("Name_%d", i),
	}

	result, err := batcher.ExecuteJob(ctx, request)

	if err != nil {
		log.Printf("Error executing job %d: %v\n", i, err)
	} else {
		log.Printf("Result of job %d: %s\n", i, result.Body)
	}

	wg.Done()
}
