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
	lambda_go "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type (
	Stack struct {
		Fn    Fn
		FnURL FnURL
		LambdaStack
	}

	LambdaStack cdk.Stack

	AppScope constructs.Construct
	Scope    constructs.Construct

	Fn    lambda_go.GoFunction
	FnURL lambda.FunctionUrl

	LogGroup logs.LogGroup
)

func NewStack(scope AppScope, environment *cdk.Environment) LambdaStack {
	return cdk.NewStack(scope, jsii.String("MicrobatchLambdaStack"), &cdk.StackProps{
		Description: jsii.String("Microbatch test stack"),
		Env:         environment,
	})
}
