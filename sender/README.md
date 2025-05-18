# PUSH Sender

`Sender`는 메시지 브로커(SNS → SQS)로부터 메시지를 받아,  
이를 전처리한 후 Session Manager로 전달하는 역할을 수행하는 마이크로서비스입니다.  
주로 다음과 같은 기능을 담당합니다.

- `SQS 메시지 수신`: SNS로부터 퍼블리시된 메세지를 SQS 큐에서 수신합니다.
- `gRPC를 통한 전달`: 수신 받은 메세지를 `Session Manager`로 전달합니다.
