package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/haandol/golang-kafka-example/pkg/config"
	"github.com/haandol/golang-kafka-example/pkg/connector/producer"
	"github.com/haandol/golang-kafka-example/pkg/util"
)

func init() {
	util.InitLogger()
}

func main() {
	logger := util.GetLogger()

	cfg := config.Load()
	logger.Infow("config", "cfg", cfg)

	p, err := producer.Connect(&cfg.Kafka)
	if err != nil {
		logger.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()

		var key int64

		for {
			select {
			case err := <-ctx.Done():
				logger.Errorw("error on job", "err", err)
				return
			default:
				now := time.Now()
				value := now.Format(time.RFC3339)

				ctx := context.Background()
				logger.Infow("produce message", "topic", "test", "key", key, "value", value)
				p.Produce(ctx, "test", fmt.Sprintf("%v", key), []byte(value))
			}

			key++
			time.Sleep(1 * time.Second)
		}
	}()

	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, os.Interrupt)

	select {
	case err := <-ctx.Done():
		logger.Errorw("error on job", "err", err)
	case <-sigs:
		cancel()
		logger.Info("User interrupt for quitting...")
	}

	ctx, cancel = context.WithTimeout(
		context.Background(),
		time.Second*time.Duration(30),
	)
	go func() {
		defer cancel()
		producer.Close(ctx)
	}()

	select {
	case err := <-ctx.Done():
		logger.Errorw("error on job", "err", err)
	case <-sigs:
		cancel()
		logger.Info("User interrupt for quitting...")
	}
}
