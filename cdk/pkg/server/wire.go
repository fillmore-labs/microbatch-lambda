//go:build wireinject
// +build wireinject

package server

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/google/wire"
)

func NewStack(scope AppScope, environment *cdk.Environment) Stack {
	wire.Build(
		wire.Value(StackProps{
			ID:          "MicrobatchLambdaStack",
			Description: "Microbatch test stack",
		}),
		newStack,
		wire.Bind(new(Scope), new(cdk.Stack)),
		NewFn,
		NewFnURL,
		wire.Struct(new(Stack), "*"),
	)

	return Stack{}
}
