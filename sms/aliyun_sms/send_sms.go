package aliyunsms

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/njmdk/common/utils"
)

type SmsStruct struct {
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	Region          string `json:"region"`
	TemplateCode    string `json:"template_code"`
	HTTPAddr        string `json:"http_addr"`
	SignName        string `json:"sign_name"`
}

type SmsSendStruct struct {
	*SmsStruct
	PhoneNumber   string
	TemplateParam string
}

func NewSmsSendStruct(base *SmsStruct) *SmsSendStruct {
	return &SmsSendStruct{SmsStruct: base}
}

func SendAliSms(sss *SmsSendStruct) error {
	query := utils.NewURLValues().
		Add("access_key_id", sss.AccessKeyID).
		Add("access_key_secret", sss.AccessKeySecret).
		Add("phone", sss.PhoneNumber).
		Add("code", sss.TemplateParam).
		Add("template_code", sss.TemplateCode).
		Add("sign_name", sss.SignName).
		Encode()
	uri := url.URL{
		Scheme:     "http",
		Host:       sss.HTTPAddr,
		Path:       "/",
		ForceQuery: false,
		RawQuery:   query,
	}
	response, err := http.Get(uri.String())
	if err != nil {
		return err
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	resp := &Response{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return err
	}
	if resp.Status == "OK" {
		return nil
	}
	return errors.New(resp.Message)
}
