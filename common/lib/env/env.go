package env

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

// 모든 서비스에서 사용하는 env.
type Env struct {
	Log     Log     `env:", prefix=LOG_"`
	App     App     `env:", prefix=APP_"`
	Aws     Aws     `env:", prefix=AWS_"`
	Linker  Linker  `env:", prefix=LINKER_"`
	Session Session `env:", prefix=SESSION_"`
	DB      DB      `env:", prefix=MARIA_"`
	Redis   Redis   `env:", prefix=REDIS_"`
	Pod     PodInfo `env:", prefix=POD_"`
}

type PodInfo struct {
	Index int `env:"INDEX"`
}

type App struct {
	Stage string `env:"STAGE, default=dev"`
	Host  string `env:"HOST, default=0.0.0.0"`
	Addr  string `env:"ADDR, default=localhost"`
}

type Log struct {
	Level string `env:"LEVEL, default=debug"`
}

type Aws struct {
	PushQueueUrl string `env:"PUSH_QUEUE_URL"`
	SnsARN       string `env:"SNS_ARN"`
}
type Linker struct {
	HttpPort string `env:"HTTP_PORT, default=8880"`
	GrpcPort string `env:"GRPC_PORT, default=50051"`
}

type Session struct {
	Port int `env:"PORT, default=50052"`
}

type DB struct {
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	Database string `env:"DATABASE"`
	Host     string `env:"HOST"`
	Port     string `env:"PORT"`
}
type Redis struct {
	Port string `env:"PORT"`
}

func validateEnv(e *Env) {
	if e.DB.Host == "" || e.DB.Database == "" || e.DB.Password == "" || e.DB.User == "" {
		log.Fatal("Invalid DB env")
	}
	if e.Aws.PushQueueUrl == "" || e.Aws.SnsARN == "" {
		log.Fatal("Ivalid AWS env")
	}
}

func NewEnv() Env {
	_ = godotenv.Load()
	env := Env{}
	if err := envconfig.Process(context.Background(), &env); err != nil {
		log.Fatal(err)
	}
	validateEnv(&env)
	return env
}
