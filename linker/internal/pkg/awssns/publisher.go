package awssns

import (
	"context"
	"encoding/json"
	"fmt"
	"push/common/lib"
	awsc "push/common/pkg/aws"
	"push/linker/internal/api/dto"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
)

type Publisher interface {
	Publish(context.Context, dto.SnsBody) (string, error)
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

func (p *publisher) Publish(ctx context.Context, msg dto.SnsBody) (string, error) {
	var messageId string

	b, _ := json.Marshal(msg)

	out, err := p.client.Publish(ctx, &sns.PublishInput{
		Message:        aws.String(string(b)),
		TopicArn:       aws.String(p.env.Aws.SnsARN),
		MessageGroupId: aws.String(fmt.Sprintf("user-%d", msg.UserId)),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"messageType": {
				DataType:    aws.String("String"),
				StringValue: aws.String("push"),
			},
		},
		MessageDeduplicationId: aws.String(fmt.Sprint(msg.MsgId)),
	})
	if err != nil {
		return messageId, err
	}
	messageId = *out.MessageId
	fmt.Println("Published message ID:", *out.MessageId)

	return messageId, nil
}
