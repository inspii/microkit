package types

import (
	"math/rand"
	"time"
)

const randBase = "0123456789abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().Unix())
}

// RandDuration 生成 [from, to) 区间内的随机时间间隔
func RandDuration(from, to time.Duration) time.Duration {
	v := rand.Intn(int(to) - int(from))
	return from + time.Duration(v)
}

// RandString 生成长度为 length 的随机字符串
func RandString(length int) string {
	buf := make([]byte, length)
	for i := 0; i < length; i++ {
		buf[i] = randBase[rand.Intn(len(randBase))]
	}
	return string(buf)
}
