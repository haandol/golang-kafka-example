package types

import (
	"context"
	"time"
)

type KafkaMessage struct {
	Topic     string
	Key       string
	Value     []byte
	Timestamp time.Time
}

type KakfaHandlerFunc func(context.Context, *KafkaMessage) error
