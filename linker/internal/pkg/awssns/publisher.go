package awssns

import (
	"context"
	"fmt"
	"push/common/lib"
	awsc "push/common/pkg/aws"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type Publisher interface {
	Publish(ctx context.Context, message string) (string, error)
}

type publisher struct {
	client *sns.Client
	env    lib.Env
}

func NewPublisher(cfg awsc.AwsConfig, env lib.Env) Publisher {
	return &publisher{
		client: sns.NewFromConfig(cfg.Config),
		env:    env,
	}
}

func (p *publisher) Publish(ctx context.Context, message string) (string, error) {
	var messageId string
	out, err := p.client.Publish(ctx, &sns.PublishInput{
		Message:        aws.String(message),
		TopicArn:       aws.String(p.env.Aws.SnsARN),
		MessageGroupId: aws.String("default"),
	})
	if err != nil {
		return messageId, err
	}
	messageId = *out.MessageId
	fmt.Println("Published message ID:", *out.MessageId)

	return messageId, nil
}
