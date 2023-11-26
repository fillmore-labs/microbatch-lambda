package server

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type (
	AppScope constructs.Construct

	StackProps struct {
		ID          string
		Description string
	}
)

func newStack(scope AppScope, props StackProps, environment *cdk.Environment) cdk.Stack {
	return cdk.NewStack(scope, jsii.String(props.ID), &cdk.StackProps{
		Description: jsii.String(props.Description),
		Env:         environment,
	})
}
