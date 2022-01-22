package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func HashMap(amap map[string]string) string {
	h := md5.New()

	for k, v := range amap {
		h.Write(StringIn(k + "=" + v + ";"))
	}
	return hex.EncodeToString(h.Sum(nil))
}
