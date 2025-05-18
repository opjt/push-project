package sqs

import (
	"context"
	"fmt"
	"log"
	"push/common/lib"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"go.uber.org/fx"
)

type Consumer struct {
	client   *sqs.Client
	queueURL string
	ctx      context.Context
	cancel   context.CancelFunc
	log      lib.Logger
}

func NewConsumer(client *sqs.Client, lc fx.Lifecycle, log lib.Logger, env lib.Env) *Consumer {
	ctx, cancel := context.WithCancel(context.Background())

	var queueURL string
	if env.App.Stage == "dev" {
		queueURL = "http://localhost:4566/000000000000/PushQueue.fifo"
	} else {
		queueURL = env.Aws.QueueUrl
	}
	c := &Consumer{
		client:   client,
		queueURL: queueURL,
		ctx:      ctx,
		cancel:   cancel,
		log:      log,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go c.start()
			return nil
		},
		OnStop: func(context.Context) error {
			log.Info("Stopping SQS Consumer...")
			c.cancel()
			return nil
		},
	})

	return c
}

func (c *Consumer) start() {
	fmt.Println("SQS Consumer started")
	for {
		select {
		case <-c.ctx.Done():
			c.log.Info("SQS Consumer stopped")
			return
		default:
			c.poll()
			time.Sleep(1 * time.Second)
		}
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
		log.Println("No messages received")
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
