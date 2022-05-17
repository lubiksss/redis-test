package main

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-redis/redis/v8"
	"log"
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
	_ = rdb.FlushDB(ctx).Err()

	//testStringType(rdb, ctx)

	testHashType(rdb, ctx)

}

func testStringType(rdb *redis.Client, ctx context.Context) {
	err := rdb.Set(ctx, "key1", "value1", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key1").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key1:", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2: ", val2)
	}
}

func testHashType(rdb *redis.Client, ctx context.Context) {
	// Set some fields.
	if _, err := rdb.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, "key", "str1", "hello")
		rdb.HSet(ctx, "key", "str2", "world")
		rdb.HSet(ctx, "key", "int", 123)
		rdb.HSet(ctx, "key", "bool", true)
		return nil
	}); err != nil {
		panic(err)
	}

	var model1, model2 Model

	// Scan all fields into the model.
	if err := rdb.HGetAll(ctx, "key").Scan(&model1); err != nil {
		panic(err)
	}

	// Or scan a subset of the fields.
	if err := rdb.HMGet(ctx, "key", "str1", "int").Scan(&model2); err != nil {
		panic(err)
	}

	log.Println(model1)
	log.Println(model2)

	spew.Dump(model1)
	spew.Dump(model2)
}

type Model struct {
	Str1 string `redis:"str1"`
	Str2 string `redis:"str2"`
	Int  int    `redis:"int"`
	Bool bool   `redis:"bool"`
}
