package web

import (
	"encoding/base64"
	"fmt"
	"github.com/shudiwsh2009/reservation_thxx_go/buslogic"
	"github.com/shudiwsh2009/reservation_thxx_go/config"
	"github.com/shudiwsh2009/reservation_thxx_go/model"
	re "github.com/shudiwsh2009/reservation_thxx_go/rerror"
	"github.com/shudiwsh2009/reservation_thxx_go/service"
	"github.com/shudiwsh2009/reservation_thxx_go/utils"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func setSession(w http.ResponseWriter, userId string, username string, userType int) error {
	now := time.Now()
	encSession, err := encryptSession(userId, username, userType, now)
	if err != nil {
		return re.NewRError("fail to encrypt session", err)
	}
	redisKeys, err := service.RedisClient().Keys(fmt.Sprintf(model.RedisKeyLogin, userType, userId, "*")).Result()
	if err != nil {
		return re.NewRError("fail to get user login session keys from redis", err)
	}
	maxLoginCount, ok := buslogic.UserLoginMaxCount[userType]
	if !ok {
		return re.NewRError(fmt.Sprintf("unknown user type: %d", userType), nil)
	}
	// remove extra redis login session
	sort.Strings(redisKeys)
	for i := 0; i <= len(redisKeys)-maxLoginCount; i++ {
		service.RedisClient().Del(redisKeys[i])
	}
	// double check redis login count
	redisKeys, err = service.RedisClient().Keys(fmt.Sprintf(model.RedisKeyLogin, userType, userId, "*")).Result()
	if err != nil {
		return re.NewRError("fail to check double redis login count", err)
	}
	if len(redisKeys) >= maxLoginCount {
		return re.NewRError("redis login count is still larger than maximum", nil)
	}
	// set current login session to redis
	maxExpire, ok := buslogic.UserLoginExpire[userType]
	if !ok {
		return re.NewRError(fmt.Sprintf("unknown user type: %d", userType), nil)
	}
	err = service.RedisClient().Set(fmt.Sprintf(model.RedisKeyLogin, userType, userId, now.Format("20060102150405")),
		now.Format("2006-01-02 15:04:05"), maxExpire).Err()
	if err != nil {
		return re.NewRError("fail to set user login session to redis", err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "SESS",
		Value:    encSession,
		Path:     "/",
		Expires:  time.Now().AddDate(1, 0, 0),
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
	})
	return nil
}

func getSession(r *http.Request) (string, string, int, error) {
	cookie, err := r.Cookie("SESS")
	if err != nil {
		return "", "", 0, re.NewRError("fail to get session from cookie", err)
	}
	encSession := cookie.Value
	userId, username, userType, loginTime, err := decryptSession(encSession)
	if err != nil {
		return "", "", 0, re.NewRError("fail to decrpt session", err)
	}
	expireDuration, ok := buslogic.UserLoginExpire[userType]
	if !ok {
		return "", "", 0, re.NewRError(fmt.Sprintf("unknown user type: %d", userType), nil)
	}
	if loginTime.Add(expireDuration).Before(time.Now()) {
		return "", "", 0, re.NewRErrorCode("session is out of date", nil, re.ErrorExpireSession)
	}
	redisLoginTime, err := service.RedisClient().Get(fmt.Sprintf(model.RedisKeyLogin, userType, userId, loginTime.Format("20060102150405"))).Result()
	if err != nil || redisLoginTime != loginTime.Format("2006-01-02 15:04:05") {
		return "", "", 0, re.NewRErrorCode("no session", err, re.ErrorNoLogin)
	}
	return userId, username, userType, nil
}

func clearSession(w http.ResponseWriter, r *http.Request) error {
	http.SetCookie(w, &http.Cookie{
		Name:   "SESS",
		Path:   "/",
		MaxAge: -1,
	})
	cookie, err := r.Cookie("SESS")
	if err != nil {
		return re.NewRError("fail to get session from cookie", err)
	}
	encSession := cookie.Value
	userId, _, userType, loginTime, err := decryptSession(encSession)
	if err != nil {
		return re.NewRError("fail to decrpt session", err)
	}
	return service.RedisClient().Del(fmt.Sprintf(model.RedisKeyLogin, userType, userId, loginTime.Format("20060102150405"))).Err()
}

func encryptSession(userId string, username string, userType int, now time.Time) (string, error) {
	key := []byte(config.Instance().SessionKeyCode)
	cookie := strings.Join([]string{userId, username, strconv.Itoa(userType), now.Format("2006-01-02 15:04:05")}, "|")
	aesCookie, err := utils.AesEncrypt([]byte(cookie), key, key)
	if err != nil {
		return "", err
	}
	base64Cookie := base64.StdEncoding.EncodeToString(aesCookie)
	return base64Cookie, nil
}

func decryptSession(encCookie string) (string, string, int, time.Time, error) {
	key := []byte(config.Instance().SessionKeyCode)
	aesCookie, err := base64.StdEncoding.DecodeString(encCookie)
	if err != nil {
		return "", "", 0, time.Now(), err
	}

	decCookie, err := utils.AesDecrypt(aesCookie, key, key)
	if err != nil {
		return "", "", 0, time.Now(), err
	}

	splits := strings.Split(string(decCookie), "|")
	if len(splits) != 4 {
		return "", "", 0, time.Now(), re.NewRError("fail to get cookie", nil)
	}

	userId := splits[0]
	username := splits[1]
	userType, err := strconv.Atoi(splits[2])
	if err != nil {
		return "", "", 0, time.Now(), err
	}
	loginTime, err := time.ParseInLocation("2006-01-02 15:04:05", splits[3], time.Local)
	if err != nil {
		return "", "", 0, time.Now(), err
	}
	return userId, username, userType, loginTime, nil
}
