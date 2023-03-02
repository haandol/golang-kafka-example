package config

import (
	"log"

	"github.com/haandol/golang-kafka-example/pkg/util"
	"github.com/joho/godotenv"
)

type App struct {
	Name  string `validate:"required"`
	Stage string `validate:"required"`
}

type Kafka struct {
	Seeds            []string `validate:"required"`
	MessageExpirySec int      `validate:"required,number"`
	BatchSize        int      `validate:"required,number"`
}

type Consumer struct {
	GroupID string `validate:"required"`
}

type Config struct {
	App      App
	Kafka    Kafka
	Consumer Consumer
}

// Load config.Config from environment variables for each stage
func Load() Config {
	stage := getEnv("APP_STAGE").String()
	log.Printf("Loading %s config\n", stage)

	if err := godotenv.Load(); err != nil {
		log.Panic("Error loading .env file")
	}

	cfg := Config{
		App: App{
			Name:  getEnv("APP_NAME").String(),
			Stage: getEnv("APP_STAGE").String(),
		},
		Kafka: Kafka{
			Seeds:            getEnv("KAFKA_SEEDS").Split(","),
			MessageExpirySec: getEnv("KAFKA_MESSAGE_EXPIRY_SEC").Int(),
			BatchSize:        getEnv("KAFKA_BATCH_SIZE").Int(),
		},
		Consumer: Consumer{
			GroupID: getEnv("CONSUMER_GROUP_ID").String(),
		},
	}

	if err := util.ValidateStruct(cfg); err != nil {
		log.Panicf("Error validating config: %s", err.Error())
	}

	return cfg
}
