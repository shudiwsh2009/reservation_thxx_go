package config

import (
	"encoding/json"
	"github.com/mijia/sweb/log"
	"io/ioutil"
)

type Config struct {
	isStaging         bool     `json:"-"`
	AppEnv            string   `json:"app_env"`
	SMSUid            string   `json:"sms_uid"`
	SMSKey            string   `json:"sms_key"`
	SMTPHost          string   `json:"smtp_host"`
	SMTPUser          string   `json:"smtp_user"`
	SMTPPassword      string   `json:"smtp_password"`
	EmailAddressAdmin []string `json:"email_address_admin"`
	EmailAddressDev   []string `json:"email_address_dev"`
	MongoHost         string   `json:"mongo_host"`
	MongoAuthDatabase string   `json:"mongo_auth_database"`
	MongoAuthUser     string   `json:"mongo_auth_user"`
	MongoAuthPassword string   `json:"mongo_auth_password"`
	MongoDatabase     string   `json:"mongo_database"`
	RedisAddress      string   `json:"redis_address"`
	RedisPassword     string   `json:"redis_password"`
	RedisDatabase     int      `json:"redis_database"`
	SessionKeyCode    string   `json:"session_key_code"`
}

var conf *Config

func (c *Config) IsStagingEnv() bool {
	if c.isStaging {
		return true
	}
	return c.AppEnv != "PRODUCTION"
}

func Init(path string) *Config {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("load config thxx.conf failed: %+v", err)
	}
	conf = &Config{}
	err = json.Unmarshal(buf, conf)
	if err != nil {
		log.Fatalf("decode config file failed: %s, err: %+v", string(buf), err)
	}
	return conf
}

func InitWithParams(path string, isStaging bool) *Config {
	conf := Init(path)
	if conf != nil {
		conf.isStaging = isStaging
	}
	return conf
}

func Instance() *Config {
	if conf == nil {
		Init("./deploy/thxx.conf")
	}
	return conf
}
