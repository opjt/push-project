package awsinfra

import (
	"context"
	"push/common/lib/env"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type AwsConfig struct {
	Config aws.Config
}

func NewAwsConfig(env env.Env) (AwsConfig, error) {
	mod := env.App.Stage

	var cfg aws.Config
	var err error
	if mod == "dev" {
		cfg, err = config.LoadDefaultConfig(context.Background(),
			config.WithEndpointResolver(aws.EndpointResolverFunc(
				func(service, region string) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:           "http://localhost:4566",
						SigningRegion: "us-east-1",
					}, nil
				})),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(context.Background())
	}

	return AwsConfig{Config: cfg}, err
}
