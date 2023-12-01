module github.com/fillmore-labs/microbatch-lambda/client

go 1.21

toolchain go1.21.4

replace github.com/fillmore-labs/microbatch-lambda/api => ../api

require (
	fillmore-labs.com/microbatch v0.0.2
	github.com/aws/aws-sdk-go-v2 v1.23.4
	github.com/aws/aws-sdk-go-v2/config v1.25.10
	github.com/fillmore-labs/microbatch-lambda/api v0.0.0-00010101000000-000000000000
	github.com/spf13/viper v1.17.0
)

require (
	github.com/aws/aws-sdk-go-v2/credentials v1.16.8 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.14.8 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.2.7 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.5.7 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.7.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.10.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.10.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.18.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.21.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.26.1 // indirect
	github.com/aws/smithy-go v1.18.1 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml/v2 v2.1.0 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20231127185646-65229373498e // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
