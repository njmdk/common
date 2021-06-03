package wx_login

import (
	"errors"
	"github.com/njmdk/common/logger"
	"github.com/njmdk/common/network/http"
	"github.com/njmdk/common/utils"
	"time"
)

type OtherLogin interface {
	GetAccessToken(string) (*WXTokenResponse, error)
	GetUserInfo(string, string) (*WXInfoResponse, error)
}

type WX struct {
	httpClient *http.Client
	log        *logger.Logger
	appId      string
	appSecret  string
}

func NewWXLogin(appId string,secretKey string,log *logger.Logger) (*WX, error) {
	c, err := http.NewHTTPClientWithTimeout("https://api.weixin.qq.com", nil, log, false, time.Second*10)
	if err != nil {
		return nil, err
	}
	c.ShowLog = false
	api := &WX{
		httpClient: c,
		log:        log,
		appId:      appId,
		appSecret:  secretKey,
	}
	return api, nil
}
type NewTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn int64 `json:"expires_in"`
	ErrCode int64 `json:"errcode"`
	ErrMsg string `json:"errmsg"`
}
func (this_ *WX) GetNewAccessToken() (*NewTokenResponse,error) {
	//https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=APPID&secret=APPSECRET
	out := &NewTokenResponse{}
	queryParams := utils.NewURLValues()
	queryParams.Add("appid", this_.appId)
	queryParams.Add("secret", this_.appSecret)
	queryParams.Add("grant_type", "client_credential")
	err := this_.httpClient.GETJsonResponse("/cgi-bin/token", queryParams, nil, nil, out)
	if err != nil {
		return nil, err
	}
	return out,nil
}
type NewQrCodeResponse struct {
	//{"ticket":"gQH47joAAAAAAAAAASxodHRwOi8vd2VpeGluLnFxLmNvbS9xL2taZ2Z3TVRtNzJXV1Brb3ZhYmJJAAIEZ23sUwMEmm
	//3sUw==","expire_seconds":60,"url":"http://weixin.qq.com/q/kZgfwMTm72WWPkovabbI"}
	Ticket string `json:"ticket"`
	ExpireSeconds int64 `json:"expire_seconds"`
	Url string `json:"url"`
	//{"errcode":48001,"errmsg":"api unauthorized rid: 5fff9e14-2c4fbbd9-52a8b525"}
	ErrCode int64 `json:"errcode"`
	ErrMsg string `json:"errmsg"`
}
type SceneInfo struct {
	SceneStr string `json:"scene_str"`
}
type ActionInfo struct {
	Scene SceneInfo `json:"scene"`
}
type NewQrCodeRequest struct {
	ExpireSeconds int64 `json:"expire_seconds"`
	ActionName string `json:"action_name"`
	ActionInfo ActionInfo `json:"action_info"`
}
func (this_ *WX) GetNewQrCode(token string,sceneStr string) (*NewQrCodeResponse,error) {
	//https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=TOKEN
	out := &NewQrCodeResponse{}
	queryParams := utils.NewURLValues()
	queryParams.Add("access_token", token)
	scene := &NewQrCodeRequest{
		ExpireSeconds: 604800,
		ActionName:    "QR_STR_SCENE",
		ActionInfo:    ActionInfo{Scene:SceneInfo{SceneStr:sceneStr}},
	}
	err := this_.httpClient.POSTJsonResponse("/cgi-bin/qrcode/create", queryParams, nil, scene, out)
	if err != nil {
		return nil, err
	}
	return out,nil
}
func (this_ *WX) GetAccessToken(code string,platform int64) (*WXTokenResponse, error) {
	out := &WXTokenResponse{}
	queryParams := utils.NewURLValues()
	switch platform {
	case 0://web
		queryParams.Add("code", code)
		queryParams.Add("appid", this_.appId)
		queryParams.Add("secret", this_.appSecret)
		queryParams.Add("grant_type", "authorization_code")
		err := this_.httpClient.GETJsonResponse("/sns/oauth2/access_token", queryParams, nil, nil, out)
		if err != nil {
			return nil, err
		}
	case 1://招聘小程序微信登录
		queryParams.Add("js_code", code)
		queryParams.Add("appid", this_.appId)
		queryParams.Add("secret", this_.appSecret)
		queryParams.Add("grant_type", "authorization_code")
		err := this_.httpClient.GETJsonResponse("/sns/jscode2session", queryParams, nil, nil, out)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("不支持的登录方式")
	}
	return out, nil
}
func (this_ *WX) GetUserInfo(accessToken string, openid string) (*WXInfoResponse, error) {
	out := &WXInfoResponse{}
	queryParams := utils.NewURLValues()
	queryParams.Add("access_token", accessToken)
	queryParams.Add("openid", openid)
	queryParams.Add("lang", "zh_CN")
	err := this_.httpClient.GETJsonResponse("/sns/userinfo", queryParams, nil, nil, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
