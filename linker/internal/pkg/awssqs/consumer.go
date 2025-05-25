package awssqs

import (
	"context"
	"log"
	"push/common/lib"
	awsc "push/common/pkg/aws"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"go.uber.org/fx"
)

type Consumer struct {
	client   *sqs.Client
	queueURL string
	ctx      context.Context
	logger   lib.Logger
}

func NewConsumer(cfg awsc.AwsConfig, lc fx.Lifecycle, logger lib.Logger, env lib.Env) *Consumer {
	client := sqs.NewFromConfig(cfg.Config)

	ctx, cancel := context.WithCancel(context.Background())

	c := &Consumer{
		client:   client,
		queueURL: env.Aws.StatusQueueUrl,
		ctx:      ctx,
		logger:   logger,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go c.start()
			return nil
		},
		OnStop: func(context.Context) error {
			logger.Info("Stopping SQS Consumer...")
			cancel()
			return nil
		},
	})

	return c
}

func (c *Consumer) start() {
	c.logger.Info("SQS Consumer started")
	for {
		c.poll()
	}

}

func (c *Consumer) poll() {
	resp, err := c.client.ReceiveMessage(c.ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &c.queueURL,
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     10,
	})
	if err != nil {
		log.Printf("ReceiveMessage failed: %v", err)
		return
	}

	if len(resp.Messages) == 0 {
		c.logger.Debug("No messages received")
		return
	}

	for _, msg := range resp.Messages {
		log.Printf("Received: %s", aws.ToString(msg.Body))

		_, err := c.client.DeleteMessage(c.ctx, &sqs.DeleteMessageInput{
			QueueUrl:      &c.queueURL,
			ReceiptHandle: msg.ReceiptHandle,
		})
		if err != nil {
			log.Printf("DeleteMessage failed: %v", err)
		} else {
			log.Println("Message deleted")
		}
	}
}
