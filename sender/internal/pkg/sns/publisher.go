package sns

import (
	"context"
	"encoding/json"
	"push/common/lib"
	"push/sender/internal/dto"

	awsc "push/common/pkg/aws"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	types "github.com/aws/aws-sdk-go-v2/service/sns/types"
)

type Publisher interface {
	Publish(context.Context, dto.Message) (string, error)
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

func (p *publisher) Publish(ctx context.Context, msg dto.Message) (string, error) {
	var messageId string

	b, _ := json.Marshal(msg)

	out, err := p.client.Publish(ctx, &sns.PublishInput{
		Message:        aws.String(string(b)),
		TopicArn:       aws.String(p.env.Aws.SnsARN),
		MessageGroupId: aws.String("default"),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"messageType": {
				DataType:    aws.String("String"),
				StringValue: aws.String("status"),
			},
		},
	})
	if err != nil {
		return messageId, err
	}
	messageId = *out.MessageId

	return messageId, nil
}
