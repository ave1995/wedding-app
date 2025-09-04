package ws

import (
	"encoding/json"
	"time"
)

type envelope[T any] struct {
	Topic string `json:"topic"`
	Data  T      `json:"data"`
	Time  string `json:"time"`
}

func wrapEvent[T any](topic string, e T) ([]byte, error) {
	return json.Marshal(envelope[T]{
		Topic: topic,
		Data:  e,
		Time:  time.Now().Format(time.RFC3339),
	})
}

func unwrapEvent[T any](msg []byte) (*envelope[T], error) {
	var env envelope[T]
	if err := json.Unmarshal(msg, &env); err != nil {
		return nil, err
	}
	return &env, nil
}
