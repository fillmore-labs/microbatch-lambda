package main

import (
	"log/slog"
	"runtime/debug"

	"github.com/aws/aws-lambda-go/lambda"
	function "github.com/fillmore-labs/microbatch-lambda/lambda"
)

func main() {
	slog.Info("Starting", "revision", revision())

	lambda.StartWithOptions(function.Handler)
}

func revision() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}

	return "development"
}
