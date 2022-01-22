package utils

import "github.com/vmihailenco/msgpack/v5"

func MsgPack(data interface{}) ([]byte, error) {
	return msgpack.Marshal(data)
}

func MsgUnpack(raw []byte, data interface{}) error {
	return msgpack.Unmarshal(raw, data)
}
