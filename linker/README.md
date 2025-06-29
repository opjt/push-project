# PUSH Linker

`Linker`는 외부 시스템 또는 내부 로직에서 생성된 알림 메시지를 수신하고, 이를 AWS SNS에 발행(Publish)하는 Publisher 역할을 수행하는 모듈입니다.
전체 아키텍처에서 메시지 흐름의 시작점이 되며, 이후 SNS → SQS → Sender 모듈로 이어지는 비동기 메시지 파이프라인을 형성합니다.

## 주요 기능

- HTTP를 통해 외부의 요청 수신 또는 내부 트리거(배치 작업, 디비 연동)을 통해 메세지를 생성합니다.
- 생성된 메세지를 AWS SNS Topic에 발행 및 DB에 데이터를 저장합니다.
