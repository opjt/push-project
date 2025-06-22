# PUSH Sender

`Sender`는 AWS SQS로부터 메시지를 소비하고, 이를 [Session Manager](../sessionmanager/README.md)에 전달하는 중간 처리자(Consumer) 역할을 수행하는 모듈입니다.  
SNS → SQS → Sender 흐름을 통해 비동기적으로 발행된 메시지를 실시간 세션 전송 시스템으로 연결합니다.

## 주요 기능

- SNS에서 퍼블리시된 메시지를 SQS 큐에서 pull 방식으로 수신합니다.
- 수신한 메시지를 전처리한 후, gRPC 호출을 통해 Session Manager에 전달합니다.
<!-- TODO(sender,README) :- 메시지 처리 실패 시 재시도 로직을 통하여 메시지 신뢰성 및 내구성을 확보합니다.  -->
