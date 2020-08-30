package redis

import (
	"context"

	"github.com/gomodule/redigo/redis"
)

type Publisher struct {
	Conn *redis.Conn
}

func (p Publisher) Publish(ctx context.Context, channel string, message interface{}) error {

	return nil
}
