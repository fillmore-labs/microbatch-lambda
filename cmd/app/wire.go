//go:build wireinject
// +build wireinject

package main

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/cxapi"
	"github.com/fillmore-labs/microbatch-lambda/pkg/server"
	"github.com/google/wire"
)

func NewAssembly() cxapi.CloudAssembly {
	wire.Build(
		CreateCloudAssembly,
		wire.Struct(new(CdkApp), "*"),
		cdk.NewApp,
		wire.Value(appProps),
		wire.Bind(new(server.AppScope), new(cdk.App)),
		env,
		server.NewStack,
	)
	return nil
}
