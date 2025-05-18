#!/bin/bash
export AWS_PAGER=""

ENDPOINT_URL="http://localhost:4566"
REGION="us-east-1"

# 1. SNS FIFO 토픽 생성
aws --endpoint-url=$ENDPOINT_URL sns create-topic \
  --name PushTopic.fifo \
  --attributes FifoTopic=true,ContentBasedDeduplication=true \
  --region $REGION || exit 1

# 2. SQS FIFO 큐 생성
aws --endpoint-url=$ENDPOINT_URL sqs create-queue \
  --queue-name PushQueue.fifo \
  --attributes FifoQueue=true,ContentBasedDeduplication=true \
  --region $REGION || exit 1

# 3. SQS 큐 URL 가져오기
QUEUE_URL=$(aws --endpoint-url=$ENDPOINT_URL sqs get-queue-url \
  --queue-name PushQueue.fifo \
  --region $REGION \
  --query 'QueueUrl' --output text) 

# 4. SQS 큐 ARN 가져오기
QUEUE_ARN=$(aws --endpoint-url=$ENDPOINT_URL sqs get-queue-attributes \
  --queue-url $QUEUE_URL \
  --attribute-names QueueArn \
  --region $REGION \
  --query 'Attributes.QueueArn' --output text)

# 5. SNS 토픽 ARN 가져오기
TOPIC_ARN=$(aws --endpoint-url=$ENDPOINT_URL sns list-topics \
  --region $REGION \
  --query "Topics[?ends_with(TopicArn, 'PushTopic.fifo')].TopicArn" \
  --output text)

# 6. SQS 정책 생성 (SNS에서 SQS로 메시지 보낼 수 있도록)
# json 형식 컨버팅해서 넣도록 개선 필요.

# 7. 정책을 SQS 큐에 설정
aws --endpoint-url=$ENDPOINT_URL sqs set-queue-attributes \
  --queue-url $QUEUE_URL \
  --attributes '{"Policy":"{\"Version\":\"2012-10-17\",\"Id\":\"sns-sqs-policy\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":\"*\",\"Action\":\"SQS:SendMessage\",\"Resource\":\"arn:aws:sqs:us-east-1:000000000000:PushQueue.fifo\",\"Condition\":{\"ArnEquals\":{\"aws:SourceArn\":\"arn:aws:sns:us-east-1:000000000000:PushTopic.fifo\"}}}]}"}'

# 8. SNS 토픽에 SQS 구독 추가
aws --endpoint-url=$ENDPOINT_URL sns subscribe \
  --topic-arn $TOPIC_ARN \
  --protocol sqs \
  --notification-endpoint $QUEUE_ARN \
  --region $REGION 

echo "SNS topic and SQS queue linked successfully!"
