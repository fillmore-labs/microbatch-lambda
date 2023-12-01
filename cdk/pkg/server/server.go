package server

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	lambda "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	lambda_go "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
)

type (
	Stack struct {
		Fn    Fn
		FnURL FnURL
		cdk.Stack
	}

	Scope constructs.Construct

	Fn    lambda_go.GoFunction
	FnURL lambda.FunctionUrl
)
