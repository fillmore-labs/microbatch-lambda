//go:build wireinject
// +build wireinject

package server

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewStack,
	wire.Bind(new(Scope), new(LambdaStack)),
	NewFn,
	NewFnURL,
	wire.Struct(new(Stack), "*"),
)
