package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
)

var (
	client = &redisClient{}
)

//redisClient struct
type redisClient struct {
	c *redis.Client
}

//GetClient get the redis client
func initialize() *redisClient {
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	if err := c.Ping().Err(); err != nil {
		panic("Unable to connect to redis " + err.Error())
	}
	client.c = c
	return client
}

//GetKey get key
func (client *redisClient) getKey(key string, src interface{}) error {
	val, err := client.c.Get(key).Result()
	if err == redis.Nil || err != nil {
		return err
	}
	err = json.Unmarshal([]byte(val), &src)
	if err != nil {
		return err
	}
	return nil
}

//GetKey get key
func (client *redisClient) getKey2(key string) (string, error) {
	val, err := client.c.Get(key).Result()
	if err == redis.Nil || err != nil {
		return val, err
	}
	return val, nil
}

//SetKey set key
func (client *redisClient) setKey(key string, value interface{}, expiration time.Duration) error {
	cacheEntry, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = client.c.Set(key, cacheEntry, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

type valueEx struct {
	Name  string
	Email string
}

func main() {
	//Use your actually ip address here
	redisClient := initialize()
	key1 := "sampleKey"
	value1 := &valueEx{Name: "someName", Email: "someemail@abc.com"}
	err := redisClient.setKey(key1, value1, time.Minute*1)
	if err != nil {
		log.Fatalf("Error: %v", err.Error())
	}
	value2 := &valueEx{}
	err = redisClient.getKey(key1, value2)
	if err != nil {
		log.Fatalf("Error: %v", err.Error())
	}
	log.Printf("Name: %s", value2.Name)
	log.Printf("Email: %s", value2.Email)

	err = redisClient.setKey("Marcelo", "OK2", time.Minute*1)
	if err != nil {
		panic(err)
	}

	valor, err := redisClient.getKey2("Marcelo")
	if err != nil {
		panic(err)
	}
	fmt.Println("key", valor)
}
