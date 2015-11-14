package sms

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/shudiwsh2009/reservation_thxx_go/util"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type SMSRequest struct {
	Uid     string `json:"Uid"`
	Key     string `json:"Key"`
	Mobile  string `json:"smsMob"`
	Content string `json:"smsText"`
}

func SendSMS(mobile string, content string) error {
	if m := util.IsMobile(mobile); !m {
		return errors.New("手机号格式不正确")
	}
	appEnv := os.Getenv("RESERVATION_THXX_ENV")
	uid := os.Getenv("RESERVATION_THXX_SMS_UID")
	key := os.Getenv("RESERVATION_THXX_SMS_KEY")
	appEnv = "ONLINE"
	uid = "shudiwsh2009"
	key = "946fee2e7ad699b065f1"
	if !strings.EqualFold(appEnv, "ONLINE") || strings.EqualFold(uid, "") || strings.EqualFold(key, "") {
		fmt.Printf("Send SMS: \"%s\" to %s.\n", content, mobile)
		return nil
	}
	requestUrl := "http://utf8.sms.webchinese.cn"
	payload := url.Values{
		"Uid":     {uid},
		"Key":     {key},
		"smsMob":  {mobile},
		"smsText": {content},
	}
	requestBody := bytes.NewBufferString(payload.Encode())
	response, err := http.Post(requestUrl, "application/x-www-form-urlencoded;charset=utf8", requestBody)
	if err != nil {
		return errors.New("短信发送失败")
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.New("短信发送失败")
	}
	fmt.Println(string(responseBody))
	if code, err := strconv.Atoi(string(response)); err != nil || code < 0 {
		return errors.New(fmt.Sprintf("短信发送失败,ErrCode:%d", code))
	}
	return nil
}
