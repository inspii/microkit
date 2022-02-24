package cache

import "github.com/go-redis/redis/v7"

// LPush 在列表中添加元素
func LPush(client *redis.Client, key string, values ...interface{}) (int64, error) {
	return client.LPush(key, values...).Result()
}

// LPop 从列表中移出一个元素
func LPop(client *redis.Client, key string) (string, error) {
	return client.LPop(key).Result()
}
