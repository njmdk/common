package aligreen

type ClientInfo struct {
	SdkVersion  string `json:"sdkVersion"`
	CfgVersion  string `json:"cfgVersion"`
	UserType    string `json:"userType"`
	UserId      string `json:"userId"`
	UserNick    string `json:"userNick"`
	Avatar      string `json:"avatar"`
	Imei        string `json:"imei"`
	Imsi        string `json:"imsi"`
	Umid        string `json:"umid"`
	Ip          string `json:"ip"`
	Os          string `json:"os"`
	Channel     string `json:"channel"`
	HostAppName string `json:"hostAppName"`
	HostPackage string `json:"hostPackage"`
	HostVersion string `json:"hostVersion"`
}

type Task struct {
	DataId  string `json:"dataId"`
	Url     string `json:"url"`
	Content string `json:"content"`
}

type BizData struct {
	BizType string   `json:"bizType"`
	Scenes  []string `json:"scenes"`
	Tasks   []*Task   `json:"tasks"`
}

type Result struct {
	Code int64 `json:"code"`
	Data []*ResultData `json:"data"`
	Msg string `json:"msg"`
}

type ResultData struct {
	Code int64 `json:"code"`
	Content string `json:"content"`
	Msg string `json:"msg"`
	Results []struct{
		Details []struct{
			Contexts struct{
				Context string `json:"context"`
			} `json:"contexts"`
			Label string `json:"label"`
		} `json:"details"`
		Label string `json:"label"`
		Rate float64 `json:"rate"`
		Scene string `json:"scene"`
		Suggestion string `json:"suggestion"`
	} `json:"results"`
}