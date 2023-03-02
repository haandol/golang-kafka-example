package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/haandol/golang-kafka-example/pkg/config"
	"github.com/haandol/golang-kafka-example/pkg/connector/consumer"
	"github.com/haandol/golang-kafka-example/pkg/types"
	"github.com/haandol/golang-kafka-example/pkg/util"
)

func init() {
	util.InitLogger()
}

func main() {
	logger := util.GetLogger()

	cfg := config.Load()
	logger.Infow("config", "cfg", cfg)

	c := consumer.NewKafkaConsumer(&cfg.Kafka, cfg.Consumer.GroupID, "test")
	c.RegisterHandler(func(ctx context.Context, msg *types.KafkaMessage) error {
		logger.Infow("message received", "group", cfg.Consumer.GroupID, "key", msg.Key, "value", string(msg.Value))
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		if err := c.Consume(ctx); err != nil {
			logger.Errorw("consumer error", "err", err)
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
		c.Close(ctx)
	}()

	select {
	case err := <-ctx.Done():
		logger.Errorw("error on job", "err", err)
	case <-sigs:
		cancel()
		logger.Info("User interrupt for quitting...")
	}
}
