package main

import (
	"os"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/cxapi"
	"github.com/aws/jsii-runtime-go"
	"github.com/fillmore-labs/microbatch-lambda/cdk/pkg/server"
)

type CdkApp struct {
	cdk.App
	serverStack server.Stack
}

func CreateCloudAssembly(app CdkApp) cxapi.CloudAssembly {
	return app.Synth(nil)
}

var appProps = &cdk.AppProps{ //nolint:gochecknoglobals
	AnalyticsReporting: jsii.Bool(false),
}

func env() *cdk.Environment {
	environment := cdk.Environment{}
	if account, ok := os.LookupEnv("CDK_DEFAULT_ACCOUNT"); ok {
		environment.Account = jsii.String(account)
	}

	if region, ok := os.LookupEnv("CDK_DEFAULT_REGION"); ok {
		environment.Region = jsii.String(region)
	}

	return &environment
}
