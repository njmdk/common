package aliyunsms

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"go.uber.org/zap"

	"github.com/njmdk/common/logger"
)

type AliYunSms struct {
	client *dysmsapi.Client
	log    *logger.Logger
}

func NewAliYunSms(accessKeyId, accessKeySecret string, log *logger.Logger) (*AliYunSms, error) {
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", accessKeyId, accessKeySecret)
	if err != nil {
		return nil, err
	}
	return &AliYunSms{
		client: client,
		log:    log,
	}, nil
}

type Response struct {
	Status  string
	Message string
	Phone   string
	Code    string
}

func (this_ *AliYunSms) SendSms(phoneNumber, signName, templateCode, code string) (out Response, err error) {
	out.Phone = phoneNumber
	out.Code = code
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.Domain = "dysmsapi.aliyuncs.com"
	request.PhoneNumbers = phoneNumber
	request.SignName = signName
	request.TemplateCode = templateCode
	request.TemplateParam = code
	var response *dysmsapi.SendSmsResponse
	response, err = this_.client.SendSms(request)
	if err != nil {
		out.Status = "sms.server.SendSms.failed"
		out.Message = err.Error()
		this_.log.Error("发送短信验证码错误", zap.Error(err), zap.String("phone", phoneNumber), zap.String("code", code), zap.Any("request", request))
		return
	}
	out.Status = response.Code
	out.Message = response.Message
	if out.Status != "OK" {
		this_.log.Error("发送短信验证码失败", zap.String("phone", phoneNumber), zap.String("code", code), zap.Any("request", request), zap.Any("response", response))
		return
	}
	this_.log.Info("发送短信验证码成功", zap.String("phone", phoneNumber), zap.String("code", code), zap.Any("request", request), zap.Any("response", response))
	return
}
