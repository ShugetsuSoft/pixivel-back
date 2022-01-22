package storage

import "hash/fnv"

func GetUrlHash(u string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(u))
	return h.Sum64()
}