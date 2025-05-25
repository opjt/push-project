package awssqs

import (
	"context"
	"errors"
	"log"
	"push/common/lib"
	awsc "push/common/pkg/aws"
	core "push/linker/internal/core/bootstrap"
	"time"

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

func NewConsumer(cfg awsc.AwsConfig, lc fx.Lifecycle, logger lib.Logger, env lib.Env, ctx *core.AppContext) *Consumer {
	client := sqs.NewFromConfig(cfg.Config)

	c := &Consumer{
		client:   client,
		queueURL: env.Aws.StatusQueueUrl,
		ctx:      ctx.Ctx,
		logger:   logger,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go c.start()
			return nil
		},
		OnStop: func(context.Context) error {
			logger.Info("Stopping SQS Consumer...")
			return nil
		},
	})

	return c
}

func (c *Consumer) start() {
	c.logger.Info("SQS Consumer started")
	for {
		select {
		case <-c.ctx.Done():
			c.logger.Info("SQS Consumer stopped")
			return
		default:
			c.poll()
		}
	}

}

func (c *Consumer) poll() {
	pollCtx, cancel := context.WithTimeout(c.ctx, 15*time.Second)
	defer cancel()

	resp, err := c.client.ReceiveMessage(pollCtx, &sqs.ReceiveMessageInput{
		// resp, err := c.client.ReceiveMessage(c.ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &c.queueURL,
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     10,
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			c.logger.Warn("SQS ReceiveMessage timed out")
		} else if errors.Is(err, context.Canceled) {
			c.logger.Info("SQS polling canceled due to shutdown")
		} else {
			c.logger.Error("SQS error:", err)
		}

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
