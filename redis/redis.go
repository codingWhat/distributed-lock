package redis

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

const (
	DelLua = `
	local client = redis.call('get', KEYS[1])
	if client == nil
	then
		return -1
	else
		if client == KEYS[2]
		then
			return redis.call('del', KEYS[1])
		else
		  return 0
		end
	end
`
)

type Locker struct {
	store    redis.Conn
	nodeName string
}

// Lock set nx ex
func (l *Locker) Lock(ctx context.Context, k string, ttl int) error {
	reply, err := redis.DoContext(l.store, ctx, "set", k, ttl, "NX", "EX", ttl)
	if err != nil {
		return err
	}
	_, err = redis.String(reply, nil)
	return err
}
func (l *Locker) UnLock(ctx context.Context, k string) error {
	_, err := redis.DoContext(l.store, ctx, "del", k)
	return err
}

// LockV1 set nx ex  + check client id (lua)
func (l *Locker) LockV1(ctx context.Context, k string, v string, ttl int) error {
	reply, err := redis.DoContext(l.store, ctx, "set", k, v, "NX", "EX", ttl)

	if err != nil {
		return err
	}
	_, err = redis.String(reply, nil)
	return err
}

func (l *Locker) UnLockV1(ctx context.Context, k string, v string) error {
	lua := redis.NewScript(1, DelLua)
	_, err := redis.Int(lua.DoContext(ctx, l.store, k, v))
	return err
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
