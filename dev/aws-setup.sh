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
  --attributes '{"Policy":"{\"Version\":\"2025-05-25\",\"Id\":\"sns-sqs-policy\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":\"*\",\"Action\":\"SQS:SendMessage\",\"Resource\":\"arn:aws:sqs:us-east-1:000000000000:PushQueue.fifo\",\"Condition\":{\"ArnEquals\":{\"aws:SourceArn\":\"arn:aws:sns:us-east-1:000000000000:PushTopic.fifo\"}}}]}"}'

# 8. SNS 토픽에 SQS 구독 추가
aws --endpoint-url=$ENDPOINT_URL sns subscribe \
  --topic-arn $TOPIC_ARN \
  --protocol sqs \
  --notification-endpoint $QUEUE_ARN \
  --region $REGION 

# 기존 큐 필터링 설정 예시 (push 메시지 전용)
SUB_ARN_PUSH=$(aws --endpoint-url=$ENDPOINT_URL sns subscribe \
  --topic-arn $TOPIC_ARN \
  --protocol sqs \
  --notification-endpoint $QUEUE_ARN \
  --region $REGION \
  --query 'SubscriptionArn' --output text)

aws --endpoint-url=$ENDPOINT_URL sns set-subscription-attributes \
  --subscription-arn $SUB_ARN_PUSH \
  --attribute-name FilterPolicy \
  --attribute-value '{"messageType": ["push"]}' \
  --region $REGION

# 새 큐 생성 및 ARN 조회 (예: StatusUpdateQueue.fifo)
aws --endpoint-url=$ENDPOINT_URL sqs create-queue \
  --queue-name StatusUpdateQueue.fifo \
  --attributes FifoQueue=true,ContentBasedDeduplication=true \
  --region $REGION || exit 1

STATUS_QUEUE_URL=$(aws --endpoint-url=$ENDPOINT_URL sqs get-queue-url \
  --queue-name StatusUpdateQueue.fifo \
  --region $REGION \
  --query 'QueueUrl' --output text)

STATUS_QUEUE_ARN=$(aws --endpoint-url=$ENDPOINT_URL sqs get-queue-attributes \
  --queue-url $STATUS_QUEUE_URL \
  --attribute-names QueueArn \
  --region $REGION \
  --query 'Attributes.QueueArn' --output text)

# 정책 설정도 꼭 해주세요 (SNS가 새 큐에 메시지 보낼 수 있도록)
aws --endpoint-url=$ENDPOINT_URL sqs set-queue-attributes \
  --queue-url $STATUS_QUEUE_URL \
  --attributes '{"Policy":"{\"Version\":\"2025-05-25\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":\"*\",\"Action\":\"SQS:SendMessage\",\"Resource\":\"'"$STATUS_QUEUE_ARN"'\",\"Condition\":{\"ArnEquals\":{\"aws:SourceArn\":\"'"$TOPIC_ARN"'\"}}}]}"}'
  

# 새 큐를 SNS 토픽에 구독 추가
SUB_ARN_STATUS=$(aws --endpoint-url=$ENDPOINT_URL sns subscribe \
  --topic-arn $TOPIC_ARN \
  --protocol sqs \
  --notification-endpoint $STATUS_QUEUE_ARN \
  --region $REGION \
  --query 'SubscriptionArn' --output text)

# 새 큐 구독에 필터 정책 적용 (status_update 메시지만 받음)
aws --endpoint-url=$ENDPOINT_URL sns set-subscription-attributes \
  --subscription-arn $SUB_ARN_STATUS \
  --attribute-name FilterPolicy \
  --attribute-value '{"messageType": ["status"]}' \
  --region $REGION

echo "SNS topic and SQS queue linked successfully!"
