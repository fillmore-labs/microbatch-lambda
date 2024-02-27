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
