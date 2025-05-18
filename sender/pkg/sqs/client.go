package sqs

import (
	"context"
	"push/common/lib"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func NewSQSClient(env lib.Env) (*sqs.Client, error) {
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

	if err != nil {
		return nil, err
	}

	return sqs.NewFromConfig(cfg), nil
}
