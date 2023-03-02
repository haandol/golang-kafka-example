package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/google/uuid"
	"github.com/haandol/golang-kafka-example/pkg/config"
	"github.com/haandol/golang-kafka-example/pkg/connector/producer"
	"github.com/haandol/golang-kafka-example/pkg/types"
	"github.com/haandol/golang-kafka-example/pkg/util"
)

const (
	topic = "test"
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
				msg := types.Message{
					Name:      "Heartbeat",
					ID:        uuid.NewString(),
					Version:   "1.0.0",
					Body:      fmt.Sprintf("tick %v", key),
					CreatedAt: time.Now().Format(time.RFC3339),
				}
				v, err := json.Marshal(msg)
				if err != nil {
					logger.Errorw("error on json marshal", "err", err)
					return
				}

				ctx := context.Background()
				logger.Infow("produce message", "topic", topic, "key", key, "msg", msg)
				p.Produce(ctx, topic, fmt.Sprintf("%v", key), v)
			}

			key++
			time.Sleep(2 * time.Second)
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
