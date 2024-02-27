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

//go:build wireinject
// +build wireinject

package main

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/cxapi"
	"github.com/fillmore-labs/microbatch-lambda/cdk/pkg/server"
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
		server.Set,
	)
	return nil
}
