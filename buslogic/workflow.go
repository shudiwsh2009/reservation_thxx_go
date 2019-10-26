package buslogic

import (
	"github.com/mijia/sweb/log"
	"github.com/shudiwsh2009/reservation_thxx_go/config"
	"github.com/shudiwsh2009/reservation_thxx_go/model"
	"gopkg.in/redis.v5"
	"time"
)

type Workflow struct {
	mongoClient *model.MongoClient
	redisClient *redis.Client
}

func NewWorkflow() *Workflow {
	var err error
	if time.Local, err = time.LoadLocation("Asia/Shanghai"); err != nil {
		log.Fatalf("初始化时区失败：%v", err)
	}
	ret := &Workflow{
		mongoClient: model.NewMongoClient(),
		redisClient: model.NewRedisClient(),
	}
	if config.Instance().StudentVerificationEnabled {
		err = ret.loadStudentVerification(config.Instance().StudentVerificationFile)
		if err != nil {
			log.Fatalf("读取学号验证文件失败：%+v", err)
		}
	}
	return ret
}

func (w *Workflow) MongoClient() *model.MongoClient {
	return w.mongoClient
}

func (w *Workflow) RedisClient() *redis.Client {
	return w.redisClient
}
