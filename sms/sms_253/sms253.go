package sms_253

import (
	"github.com/guonaihong/gout"
)

type Struct253SMSInfo struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Msg      string `json:"msg"`
	Phone    string `json:"phone"`
	SendTime string `json:"sendtime"`
	Report   string `json:"report"`
	Extend   string `json:"extend"`
	UID      string `json:"uid"`
}

type Struct253SMS struct {
	Addr string
	Info *Struct253SMSInfo
}

func NewStruct253SMS() *Struct253SMS {
	return &Struct253SMS{
		Addr: "http://smssh1.253.com/msg/send/json",
		Info: &Struct253SMSInfo{
			Account:  "",
			Password: "",
		},
	}
}

type Response struct {
	Code     string `json:"code"`
	MsgID    string `json:"msgId"`
	ErrorMsg string `json:"errorMsg"`
	Time     string `json:"time"`
}

func Post253SMS(s *Struct253SMS) (*Response, error) {
	out := &Response{}
	err := gout.POST(s.Addr).
		SetHeader(gout.H{"Content-Type": "application/json; charset=UTF-8"}).
		SetJSON(s.Info).
		BindJSON(out).
		Do()
	if err != nil {
		return nil, err
	}
	return out, nil
}
