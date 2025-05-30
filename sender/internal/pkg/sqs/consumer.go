package sqs

import (
	"context"
	"push/common/lib"
	awsc "push/common/pkg/aws"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"go.uber.org/fx"
)

type Consumer struct {
	client    *sqs.Client
	queueURL  string
	ctx       context.Context
	log       lib.Logger
	msgChan   chan types.Message
	workerNum int
	handler   Handler
}

func NewConsumer(cfg awsc.AwsConfig, lc fx.Lifecycle, log lib.Logger, env lib.Env, handler Handler) *Consumer {
	ctx, cancel := context.WithCancel(context.Background())
	client := sqs.NewFromConfig(cfg.Config)

	c := &Consumer{
		client:    client,
		queueURL:  env.Aws.PushQueueUrl,
		ctx:       ctx,
		log:       log,
		msgChan:   make(chan types.Message, 100),
		workerNum: 5,
		handler:   handler,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go c.start()
			return nil
		},
		OnStop: func(context.Context) error {
			log.Info("Stopping SQS Consumer...")
			cancel()
			close(c.msgChan)
			return nil
		},
	})

	return c
}

func (c *Consumer) start() {
	c.log.Info("SQS Consumer started")

	// 워커 시작
	for i := 0; i < c.workerNum; i++ {
		go c.worker(i)
	}

	for {
		select {
		case <-c.ctx.Done():
			c.log.Info("SQS Consumer stopped")
			return
		default:
			c.poll()
		}
	}
}

func (c *Consumer) poll() {
	messages, err := c.receiveMessages()
	if err != nil {
		c.log.Errorf("ReceiveMessage failed: %v", err)
		return
	}

	if len(messages) == 0 {
		c.log.Info("No messages received")
		return
	}

	for _, msg := range messages {
		select {
		case c.msgChan <- msg:
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Consumer) receiveMessages() ([]types.Message, error) {
	resp, err := c.client.ReceiveMessage(c.ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &c.queueURL,
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     10,
	})
	if err != nil {
		return nil, err
	}
	return resp.Messages, nil
}

func (c *Consumer) worker(id int) {
	c.log.Debugf("Worker %d started", id)
	for msg := range c.msgChan {
		c.processMessage(msg)
	}
}

func (c *Consumer) processMessage(msg types.Message) {
	err := c.handler.HandleMessage(c.ctx, msg)
	if err != nil {
		c.log.Errorf("Message handling failed: %v", err)
		// 선택적으로 재처리 로직 수행
		return
	}
	if err := c.deleteMessage(msg); err != nil {
		c.log.Errorf("Failed to delete message: %v", err)
	} else {
		c.log.Info("Message deleted")
	}
}

func (c *Consumer) deleteMessage(msg types.Message) error {
	_, err := c.client.DeleteMessage(c.ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &c.queueURL,
		ReceiptHandle: msg.ReceiptHandle,
	})
	return err
}
