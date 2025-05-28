package lib

import (
	"context"
	"log"
	"sync"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Env struct {
	Log    LogConfig `env:", prefix=LOG_"`
	App    App       `env:", prefix=APP_"`
	Aws    Aws       `env:", prefix=AWS_"`
	Linker Linker    `env:", prefix=LINKER_"`
	DB     DB        `env:", prefix=MARIA_"`
}

type App struct {
	Stage string `env:"STAGE, default=dev"`
}

type LogConfig struct {
	Level string `env:"LEVEL, default=debug"`
}

type Aws struct {
	PushQueueUrl string `env:"PUSH_QUEUE_URL"`
	SnsARN       string `env:"SNS_ARN"`
}
type Linker struct {
	Port string `env:"PORT, default=8880"`
}

type DB struct {
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	Database string `env:"DATABASE"`
	Host     string `env:"HOST"`
	Port     string `env:"PORT"`
}

var (
	once sync.Once
	env  Env
)

func LoadEnv() Env {
	_ = godotenv.Load()
	if err := envconfig.Process(context.Background(), &env); err != nil {
		log.Fatal(err)
	}

	return env
}

func NewEnv() Env {

	once.Do(func() {
		env = LoadEnv()

	})
	return env
}
