package cache

import (
	"strconv"

	"github.com/go-redis/redis/v7"
)

// SAdd 加入数据到集合
func SAdd(client *redis.Client, key string, values ...string) error {
	if len(values) == 0 {
		return nil
	}

	var args []interface{}
	for _, r := range values {
		args = append(args, r)
	}
	if _, err := client.SAdd(key, args...).Result(); err != nil {
		return err
	}
	return nil
}

// SRem 从集合移除数据
func SRem(client *redis.Client, key string, values ...string) error {
	if len(values) == 0 {
		return nil
	}

	var args []interface{}
	for _, r := range values {
		args = append(args, r)
	}
	if _, err := client.SRem(key, args...).Result(); err != nil && err != redis.Nil {
		return err
	}
	return nil
}

// SMembers 返回集合中所有成员
func SMembers(client *redis.Client, key string) ([]string, error) {
	members, err := client.SMembers(key).Result() // key不存在，也不回返回redis.Nil
	if err != nil {
		return nil, err
	}
	return members, nil
}

// SAddInt 加入Int类型数据到集合
func SAddInt(client *redis.Client, key string, values ...int) error {
	if len(values) == 0 {
		return nil
	}

	var args []interface{}
	for _, r := range values {
		args = append(args, r)
	}
	if _, err := client.SAdd(key, args...).Result(); err != nil {
		return err
	}
	return nil
}

// SRemInt 从集合中移除Int类型数据
func SRemInt(client *redis.Client, key string, values ...int) error {
	if len(values) == 0 {
		return nil
	}

	var args []interface{}
	for _, r := range values {
		args = append(args, r)
	}
	if _, err := client.SRem(key, args...).Result(); err != nil && err != redis.Nil {
		return err
	}
	return nil
}

// SMembersInt 返回集合中所有Int类型的成员
func SMembersInt(client *redis.Client, key string) ([]int, error) {
	members, err := client.SMembers(key).Result() // key不存在，也不回返回redis.Nil
	if err != nil {
		return nil, err
	}

	var nums []int
	for _, v := range members {
		num, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		nums = append(nums, num)
	}

	return nums, nil
}
