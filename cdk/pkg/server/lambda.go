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

package server

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	lambda "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	logs "github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	lambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/jsii-runtime-go"
)

func NewLogGroup(scope Scope) LogGroup {
	return logs.NewLogGroup(scope, jsii.String("LogGroup"), &logs.LogGroupProps{
		LogGroupClass: logs.LogGroupClass_STANDARD,
		RemovalPolicy: cdk.RemovalPolicy_DESTROY,
		Retention:     logs.RetentionDays_THREE_DAYS,
	})
}

func NewFn(scope Scope, logGroup LogGroup) Fn {
	fn := lambdago.NewGoFunction(scope, jsii.String("Handler"), &lambdago.GoFunctionProps{
		Entry:     jsii.String("lambda/cmd/bootstrap"),
		ModuleDir: jsii.String("lambda"),
		Bundling: &lambdago.BundlingOptions{
			GoBuildFlags: jsii.Strings("-ldflags \"-s -w\""),
			GoProxies: &[]*string{
				lambdago.GoFunction_GOOGLE_GOPROXY(),
				jsii.String("direct"),
			},
		},

		Runtime:      lambda.Runtime_PROVIDED_AL2023(),
		Architecture: lambda.Architecture_ARM_64(),

		LogGroup: logGroup,
	})

	return fn
}

func NewFnURL(scope Scope, fn Fn) FnURL {
	fnURL := fn.AddFunctionUrl(&lambda.FunctionUrlOptions{
		AuthType: lambda.FunctionUrlAuthType_AWS_IAM,
	})

	cdk.NewCfnOutput(scope, jsii.String("ServerUrl"), &cdk.CfnOutputProps{
		Value: fnURL.Url(),
	})

	cdk.NewCfnOutput(scope, jsii.String("ServerRegion"), &cdk.CfnOutputProps{
		Value: fnURL.Env().Region,
	})

	return fnURL
}
