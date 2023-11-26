package server

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	lambda "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	logs "github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	lambda_go "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/jsii-runtime-go"
)

func NewFn(scope Scope) Fn {
	fn := lambda_go.NewGoFunction(scope, jsii.String("ServerFunction"), &lambda_go.GoFunctionProps{
		Entry: jsii.String("./cmd/bootstrap"),

		Runtime:      lambda.Runtime_PROVIDED_AL2023(),
		Architecture: lambda.Architecture_ARM_64(),

		LogRetention: logs.RetentionDays_THREE_DAYS,
	})

	return fn
}

func NewFnURL(scope Scope, fn Fn) FnURL {
	fnURL := fn.AddFunctionUrl(&lambda.FunctionUrlOptions{
		AuthType: lambda.FunctionUrlAuthType_AWS_IAM,
	})

	cdk.NewCfnOutput(scope, jsii.String("ServerUrl"), &cdk.CfnOutputProps{
		Value:      fnURL.Url(),
		ExportName: jsii.String("ServerUrl"),
	})

	cdk.NewCfnOutput(scope, jsii.String("ServerRegion"), &cdk.CfnOutputProps{
		Value:      fnURL.Env().Region,
		ExportName: jsii.String("ServerRegion"),
	})

	return fnURL
}
