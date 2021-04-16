package sms_ps

import (
	"encoding/base64"
	"fmt"

	"github.com/guonaihong/gout"
)

type Message struct {
	To     string `json:"to"`
	Status struct {
		ID          int64  `json:"id"`
		GroupID     int64  `json:"groupId"`
		GroupName   string `json:"groupName"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"status"`
	SmsCount  int64  `json:"smsCount"`
	MessageID string `json:"messageId"`
}

type Response struct {
	Messages []Message `json:"messages"`
}

func Post(phone string, code string) (*Response, error) {
	resp := &Response{}
	err := gout.New().POST("https://api.infobip.com/sms/1/text/single").
		SetHeader(gout.H{"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte("user:xxxxxxxxxx")),
			"Accept": "application/json"}).
		SetJSON(gout.H{"from": "baiying", "to": "+86" + phone, "text": fmt.Sprintf("【测试】您的验证码是: %s ", code)}).BindJSON(resp).Do()
	if err == nil {
		return resp, nil
	}
	return nil, err
}

func GetResult() error {
	err := gout.New().GET("https://api.infobip.com/sms/1/reports").Debug(true).
		SetHeader(gout.H{"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte("user:xxxxxxxxxx")),
			"Accept": "application/json"}).
		SetQuery(gout.H{"messageId": "111111111111"}).Do()
	return err
}
