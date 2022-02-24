package cache

import (
	"time"

	"github.com/go-redis/redis/v7"
)

// HDel 删除哈希表的键值对
func HDel(client *redis.Client, key string, fields ...string) error {
	return client.HDel(key, fields...).Err()
}

// HMGet 获取哈希表的所有键值对
func HMGet(client *redis.Client, key string) (map[string]string, error) {
	return client.HGetAll(key).Result()
}

// HMSet 同时将多个域值对设置到哈希表中。
func HMSet(client *redis.Client, key string, data map[string]interface{}, expiration time.Duration) error {
	if len(data) == 0 {
		return nil
	}

	if err := client.HMSet(key, data).Err(); err != nil {
		return err
	}
	if err := client.Expire(key, expiration).Err(); err != nil {
		return err
	}
	return nil
}
