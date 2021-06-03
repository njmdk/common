package wx_login

import "encoding/json"

type WXTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionId      string `json:"unionid"`
	SessionKey string `json:"session_key"`

	ErrCode int64  `json:"errcode"` //:40029,
	ErrMsg  string `json:"errmsg"`  //:"invalid code, hints: [ req_id: RedD4fMre-qogjAA ]"
}

type WXInfoResponse struct {
	Openid     string          `json:"openid"`
	Nickname   string          `json:"nickname"`
	Sex        int64           `json:"sex"`
	Language   string          `json:"language"`
	City       string          `json:"city"`
	Province   string          `json:"province"`
	Country    string          `json:"country"`
	HeadImgUrl string          `json:"headimgurl"`
	Privilege  json.RawMessage `json:"privilege"`
	UnionId    string          `json:"unionid"`

	ErrCode int64  `json:"errcode"` //:40029,
	ErrMsg  string `json:"errmsg"`  //:"invalid code, hints: [ req_id: RedD4fMre-qogjAA ]"
}
