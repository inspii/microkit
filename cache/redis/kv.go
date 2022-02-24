package cache

import (
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"github.com/go-redis/redis/v7"
)

// Exists Key是否存在
func Exists(client *redis.Client, key string) (bool, error) {
	num, err := client.Exists(key).Result()
	return num == 1, err
}

// GetJSON 获取JSON格式数据
func GetJSON(client *redis.Client, key string, v interface{}) error {
	value, err := client.Get(key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(value), v)
}

// SetJSON 设置JSON格式数据
func SetJSON(client *redis.Client, key string, v interface{}, expiration time.Duration) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return client.Set(key, string(bytes), expiration).Err()
}

// MGetJSON 批量获取JSON格式数据
func MGetJSON(client *redis.Client, keys []string, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("mgetJSON: v must be ptr to slice")
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Slice {
		return errors.New("mgetJSON: v must be ptr to slice")
	}
	rt := rv.Type()

	values, err := client.MGet(keys...).Result()
	if err != nil {
		return err
	}

	elems := make([]reflect.Value, 0)
	for _, v := range values {
		et := rt.Elem()
		elem := reflect.New(et)
		if val, ok := v.(string); ok {
			if err := json.Unmarshal([]byte(val), elem.Interface()); err != nil {
				return err
			}
			elems = append(elems, elem.Elem())
		}
	}
	arr := reflect.Append(rv, elems...)
	rv.Set(arr)
	return nil
}

// MSetJSON 设置JSON格式数据
func MSetJSON(client *redis.Client, slice interface{}, getKey func(i int) string, expiration time.Duration) error {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return errors.New("mgetJSON: v must be ptr to slice")
	}
	if rv.Len() == 0 {
		return nil
	}

	data := make(map[string]interface{})
	for i := 0; i < rv.Len(); i++ {
		key := getKey(i)
		v := rv.Index(i).Interface()

		bytes, err := json.Marshal(v)
		if err != nil {
			return err
		}
		data[key] = string(bytes)
	}
	if _, err := client.MSet(data).Result(); err != nil {
		return err
	}
	for k := range data {
		client.Expire(k, expiration)
	}
	return nil
}

// Delete 删除数据
func Delete(client *redis.Client, key string) error {
	if _, err := client.Del(key).Result(); err != nil && err != redis.Nil {
		return err
	}
	return nil
}

func Expire(client *redis.Client, key string, expiration time.Duration) error {
	if _, err := client.Expire(key, expiration).Result(); err != nil {
		return err
	}
	return nil
}
