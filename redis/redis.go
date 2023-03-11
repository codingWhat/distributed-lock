package redis

import (
	"context"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

type Locker struct {
	store redis.Conn
}

func (l *Locker) Lock(ctx context.Context, k string, ttl int) error {
	return nil
}

func (l *Locker) UnLock(ctx context.Context, k string) error {
	return nil
}

func NewLocker() (*Locker, error) {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return nil, err
	}

	return &Locker{
		store: c,
	}, nil
}
