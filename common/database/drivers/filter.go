package drivers

import redisbloom "github.com/RedisBloom/redisbloom-go"

type CuckooFilter struct {
	cli *redisbloom.Client
}

func (red *RedisPool) NewCuckooFilter(name string) *CuckooFilter {
	client := redisbloom.NewClientFromPool(red.RedisPool, name)
	return &CuckooFilter{
		cli: client,
	}
}

func NewCuckooFilter(pool *RedisPool, name string) *CuckooFilter {
	client := redisbloom.NewClientFromPool(pool.RedisPool, name)
	return &CuckooFilter{
		cli: client,
	}
}

func (ckft *CuckooFilter) Add(name string, key string) (bool, error) {
	return ckft.cli.CfAdd(name, key)
}

func (ckft *CuckooFilter) Exists(name string, key string) (bool, error) {
	return ckft.cli.CfExists(name, key)
}

func (ckft *CuckooFilter) Del(name string, key string) (bool, error) {
	return ckft.cli.CfDel(name, key)
}
