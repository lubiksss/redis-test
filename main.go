package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
)

func main() {
	addr := os.Getenv("TEST_REDIS_ADDR")
	password := os.Getenv("TEST_REDIS_PASS")
	//log.Println(addr, password)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	var ctx = context.Background()

	err := rdb.Set(ctx, "jake", "test", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "jake").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("jake", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}
