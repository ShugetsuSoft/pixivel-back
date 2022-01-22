package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"unsafe"
)

func HashStruct(item interface{}) string {
	jsonBytes, _ := json.Marshal(item)
	return fmt.Sprintf("%x", md5.Sum(jsonBytes))
}

func StringOut(bye []byte) string {
	return *(*string)(unsafe.Pointer(&bye))
}

func StringIn(strings string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&strings))
	return *(*[]byte)(unsafe.Pointer(&[3]uintptr{x[0], x[1], x[1]}))
}

func IntIn(n int) []byte {
	data := int64(n)
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
}

func IntOut(bye []byte) int {
	bytebuff := bytes.NewBuffer(bye)
	var data int64
	binary.Read(bytebuff, binary.BigEndian, &data)
	return int(data)
}

func UintIn(n uint) []byte {
	data := uint64(n)
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
}

func UintOut(bye []byte) uint {
	bytebuff := bytes.NewBuffer(bye)
	var data uint64
	binary.Read(bytebuff, binary.BigEndian, &data)
	return uint(data)
}

func Itoa(a interface{}) string {
	if v, p := a.(int); p {
		return strconv.Itoa(v)
	}
	if v, p := a.(int16); p {
		return strconv.Itoa(int(v))
	}
	if v, p := a.(int32); p {
		return strconv.Itoa(int(v))
	}
	if v, p := a.(uint); p {
		return strconv.Itoa(int(v))
	}
	if v, p := a.(uint64); p {
		return strconv.Itoa(int(v))
	}
	if v, p := a.(float32); p {
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	}
	if v, p := a.(float64); p {
		return strconv.FormatFloat(v, 'f', -1, 32)
	}
	return ""
}

func Atoi(s string) uint64 {
	i, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0
	}
	return i
}

func Atois(s []string) []uint64 {
	res := make([]uint64, len(s))
	for i, str := range s {
		res[i] = Atoi(str)
	}
	return res
}
