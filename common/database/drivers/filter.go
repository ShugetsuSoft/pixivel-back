package drivers

import redisbloom "github.com/RedisBloom/redisbloom-go"

type BloomFilter struct {
	cli *redisbloom.Client
}

func (red *RedisPool) NewBloomFilter(name string) *BloomFilter {
	client := redisbloom.NewClientFromPool(red.RedisPool, name)
	return &BloomFilter{
		cli: client,
	}
}

func NewBloomFilter(pool *RedisPool, name string) *BloomFilter {
	client := redisbloom.NewClientFromPool(pool.RedisPool, name)
	return &BloomFilter{
		cli: client,
	}
}

func (ckft *BloomFilter) Add(name string, key string) (bool, error) {
	return ckft.cli.Add(name, key)
}

func (ckft *BloomFilter) Exists(name string, key string) (bool, error) {
	return ckft.cli.Exists(name, key)
}
