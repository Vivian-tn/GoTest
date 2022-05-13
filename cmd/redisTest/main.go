package main

import (
	"context"
	"fmt"
	"testing"

	"gopkg.in/redis.v4"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func main() {
	var t *testing.T
	TestString(t)
}
func TestString(t *testing.T) {
	rdb.Set(ctx, "key", "value", 0)
	result, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", result)
}

func TestSet(t *testing.T) {
	addCmd := rdb.SAdd(ctx, "set", "s1", "s2", "s3")
	fmt.Println(addCmd)
	stringSliceCmd := rdb.SMembers(ctx, "set")
	for _, v := range stringSliceCmd.Val() {
		fmt.Println(v)
	}
}

func TestZSet(t *testing.T) {
	intCmd := rdb.ZAdd(ctx, "zset",
		&redis.Z{Score: 10, Member: "GetcharZp"},
		&redis.Z{Score: 50, Member: "B-zero"},
		&redis.Z{Score: 20, Member: "GetcharMmc"})
	fmt.Println(intCmd)
	zRange := rdb.ZRange(ctx, "zset", 0, 10)
	for _, v := range zRange.Val() {
		fmt.Println(v)
	}
}

func TestHash(t *testing.T) {
	intCmd := rdb.HSet(ctx, "hash", map[string]string{"k1": "v1", "k2": "v2", "k3": "v3"})
	fmt.Println(intCmd)
	// 通过指定map中的key获取值
	getOne := rdb.HGet(ctx, "hash", "k1")
	fmt.Println(getOne.Val())
	all := rdb.HGetAll(ctx, "hash")
	for key, value := range all.Val() {
		fmt.Println("key --> ", key, " value --> ", value)
	}
}

func TestList(t *testing.T) {
	intCmd := rdb.LPush(ctx, "list", "l1", "l2", "l3")
	fmt.Println(intCmd)
	lRange := rdb.LRange(ctx, "list", 0, 3) // 从最左边开始取数据
	for _, v := range lRange.Val() {
		fmt.Println(v)
	}
}
