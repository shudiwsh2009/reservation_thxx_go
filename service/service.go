package service

import (
	"github.com/mijia/sweb/log"
	"github.com/shudiwsh2009/reservation_thxx_go/buslogic"
	"github.com/shudiwsh2009/reservation_thxx_go/config"
	"github.com/shudiwsh2009/reservation_thxx_go/model"
	"gopkg.in/redis.v5"
)

var wf *buslogic.Workflow

func InitService(confPath string, isStaging bool) {
	config.InitWithParams(confPath, isStaging)
	log.Infof("config loaded: %+v", *config.Instance())
	wf = buslogic.NewWorkflow()
}

func Workflow() *buslogic.Workflow {
	return wf
}

func MongoClient() *model.MongoClient {
	return wf.MongoClient()
}

func RedisClient() *redis.Client {
	return wf.RedisClient()
}
