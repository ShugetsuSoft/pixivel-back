package drivers

import (
	"context"
	"github.com/ShugetsuSoft/pixivel-back/common/utils"
	"time"

	"github.com/gomodule/redigo/redis"
)

type RedisPool struct {
	RedisPool *redis.Pool
}

func NewRedisPool(redisURL string) *RedisPool {
	return &RedisPool{
		RedisPool: &redis.Pool{
			MaxIdle:     100,
			MaxActive:   0,
			Wait:        true,
			IdleTimeout: time.Duration(100) * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.DialURL(redisURL)
				if err != nil {
					return nil, err
				}
				return c, err
			},
		},
	}
}

type RedisClient struct {
	conn redis.Conn
}

func (red *RedisPool) NewRedisClient() *RedisClient {
	conn := red.RedisPool.Get()
	return &RedisClient{
		conn: conn,
	}
}

func NewRedisClient(redisURL string) (*RedisClient, error) {
	conn, err := redis.DialURL(redisURL)
	if err != nil {
		return nil, err
	}
	return &RedisClient{
		conn: conn,
	}, nil
}

func (red *RedisClient) IsExist(err error) bool {
	return err == nil
}

func (red *RedisClient) IsError(err error) bool {
	return err != redis.ErrNil
}

func (red *RedisClient) GetValue(key string) (string, error) {
	value, err := redis.String(red.conn.Do("GET", key))
	return value, err
}

func (red *RedisClient) GetValueBytes(key string) ([]byte, error) {
	value, err := redis.Bytes(red.conn.Do("GET", key))
	return value, err
}

func (red *RedisClient) SetValue(key string, value string) error {
	_, err := red.conn.Do("SET", key, value)
	return err
}

func (red *RedisClient) SetValueBytes(key string, value []byte) error {
	_, err := red.conn.Do("SET", key, value)
	return err
}

func (red *RedisClient) DeleteValue(key string) error {
	_, err := red.conn.Do("DEL", key)
	return err
}

func (red *RedisClient) SetValueExpire(key string, value string, exp int) error {
	_, err := red.conn.Do("SETEX", key, exp, value)
	return err
}

func (red *RedisClient) SetValueBytesExpire(key string, value []byte, exp int) error {
	_, err := red.conn.Do("SETEX", key, exp, value)
	return err
}

func (red *RedisClient) ClearDB() error {
	_, err := red.conn.Do("FLUSHDB")
	return err
}

func (red *RedisClient) SetExpire(key string, exp string) error {
	_, err := red.conn.Do("EXPIRE", key, exp)
	return err
}
func (red *RedisClient) GetExpire(key string) (string, error) {
	value, err := redis.String(red.conn.Do("TTL", key))
	if err != nil {
		if err == redis.ErrNil {
			return "nil", err
		}
		return "", err
	}
	return value, nil
}

func (red *RedisClient) KeyExist(key string) (bool, error) {
	exist, err := redis.Bool(red.conn.Do("EXISTS", key))
	if err != nil {
		return false, err
	} else {
		return exist, nil
	}
}

func (red *RedisClient) Increase(key string) (int, error) {
	return redis.Int(red.conn.Do("INCR", key))
}

func (red *RedisClient) IncreaseBy(key string, num uint) (int, error) {
	return redis.Int(red.conn.Do("INCRBY", key, num))
}

func (red *RedisClient) Decrease(key string) (int, error) {
	return redis.Int(red.conn.Do("DECR", key))
}

func (red *RedisClient) DecreaseBy(key string, num uint) (int, error) {
	return redis.Int(red.conn.Do("DECRBY", key, num))
}

func (red *RedisClient) Publish(key string, message string) error {
	_, err := red.conn.Do("PUBLISH", key, message)
	return err
}

func (red *RedisClient) Subscribe(ctx context.Context, key string) (chan string, error) {
	psc := redis.PubSubConn{Conn: red.conn}
	if err := psc.PSubscribe(key); err != nil {
		return nil, err
	}
	msg := make(chan string, 5)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(msg)
				return
			default:
				switch v := psc.Receive().(type) {
				case redis.Message:
					msg <- utils.StringOut(v.Data)
				}
			}
		}
	}()
	return msg, nil
}

func (red *RedisClient) Close() {
	red.conn.Close()
}
