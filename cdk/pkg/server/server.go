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
