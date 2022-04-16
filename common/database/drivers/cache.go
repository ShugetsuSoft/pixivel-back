package drivers

import (
	"bytes"
	"github.com/ShugetsuSoft/pixivel-back/common/utils/telemetry"
	"github.com/gomodule/redigo/redis"
	"github.com/vmihailenco/msgpack/v5"
)

type Cache struct {
	redis *RedisPool
}

func NewCache(redis *RedisPool) *Cache {
	return &Cache{
		redis: redis,
	}
}

func marshal(v interface{}) ([]byte, error) {
	enc := msgpack.NewEncoder(nil)
	var buf bytes.Buffer
	enc.Reset(&buf)
	enc.SetCustomStructTag("json")

	err := enc.Encode(v)
	b := buf.Bytes()

	if err != nil {
		return nil, err
	}
	return b, err
}

func unmarshal(data []byte, v interface{}) error {
	dec := msgpack.NewDecoder(nil)
	dec.Reset(bytes.NewReader(data))
	dec.SetCustomStructTag("json")

	err := dec.Decode(v)

	return err
}

func calcKey(api string, params []string) string {
	key := api + "-"
	for k := range params {
		key += params[k] + "."
	}
	return key
}

func (c *Cache) Get(api string, params ...string) (interface{}, error) {
	key := calcKey(api, params)
	cli := c.redis.NewRedisClient()
	defer cli.Close()
	cached, err := cli.GetValueBytes(key)
	if err != nil {
		if err == redis.ErrNil {
			return nil, nil
		}
		return nil, err
	}
	var result interface{}
	err = unmarshal(cached, &result)
	telemetry.RequestsHitCache.Inc()
	return result, err
}

func (c *Cache) Set(api string, value interface{}, expire int, params ...string) error {
	key := calcKey(api, params)
	cli := c.redis.NewRedisClient()
	defer cli.Close()
	binval, err := marshal(value)
	if err != nil {
		return err
	}
	return cli.SetValueBytesExpire(key, binval, expire)
}

func (c *Cache) Clear(api string, params ...string) error {
	key := calcKey(api, params)
	cli := c.redis.NewRedisClient()
	defer cli.Close()
	err := cli.DeleteValue(key)
	if err == redis.ErrNil {
		return nil
	}
	return err
}
