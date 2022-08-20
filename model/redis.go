package model

import (
	"github.com/shudiwsh2009/reservation_thxx_go/config"
	"gopkg.in/redis.v5"
	"log"
)

const (
	RedisKeyLogin = "thxxfzzx#user_login_%d_%s_%s"
)

func NewRedisClient() *redis.Client {
	var client *redis.Client
	if config.Instance().IsStagingEnv() {
		client = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
	} else if config.Instance().IsTsinghuaEnv() {
		client = redis.NewClient(&redis.Options{
			Addr:     config.Instance().RedisAddress,
			Password: "",
			DB:       0,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:     config.Instance().RedisAddress,
			Password: config.Instance().RedisPassword,
			DB:       config.Instance().RedisDatabase,
		})
	}
	pong, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("连接Redis失败：%v", err)
	}
	log.Printf("连接Redis成功：%s", pong)
	return client
}
