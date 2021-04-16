package aligreen

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
	
	"go.uber.org/zap"

	"github.com/njmdk/common/logger"
)

type DefaultClient struct {
	Profile *Profile
	client*http.Client
}

var defaultClient = &DefaultClient{
	Profile: &Profile{AccessKeyId: "", AccessKeySecret:""},
	client:  &http.Client{
		Transport: &http.Transport{
			Proxy: nil,
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				return net.DialTimeout(network, addr, time.Second*6)
			},
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 1000,
			IdleConnTimeout:     time.Minute * 5,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: time.Second*6,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 100 {
				return errors.New("stopped after 100 redirects")
			}
			return nil
		},
	},
}

func (this_* DefaultClient) GetResponse(path string, clientInfo *ClientInfo, bizData *BizData,log*logger.Logger) []byte{
	clientInfoJson, _ := json.Marshal(clientInfo)
	bizDataJson, _ := json.Marshal(bizData)

	
	req, err := http.NewRequest(method, host + path + "?clientInfo=" + url.QueryEscape(string(clientInfoJson)), strings.NewReader(string(bizDataJson)))

	if err != nil {
		return ErrorResult(err)
	}
	addRequestHeader(string(bizDataJson), req, string(clientInfoJson), path, this_.Profile.AccessKeyId, this_.Profile.AccessKeySecret)

	response, err := this_.client.Do(req)
	if err!=nil{
		return ErrorResult(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ErrorResult(err)
	}
	return body
}

func (this_* DefaultClient)CheckText(dataID string,text string,log*logger.Logger)(success bool,err error)  {
	path := "/green/text/scan"
	
	clientInfo := &ClientInfo{}
	
	// request data
	bizType := "Green"
	scenes := []string{"antispam"}
	
	task := &Task{DataId: dataID, Content:text}
	tasks := []*Task{task}
	
	bizData := &BizData{bizType, scenes, tasks}
	resp:=this_.GetResponse(path, clientInfo, bizData,log)
	defer func() {
		if !success || err!=nil{
			if log!=nil{
				log.Error("aligreen failed",zap.Error(err),zap.ByteString("resp",resp))
			}
		}
	}()
	r:=&Result{}
	err =json.Unmarshal(resp,r)
	if err!=nil{
		return false,err
	}
	if r.Code!=200{
		return false,nil
	}
	if len(r.Data)!=1 {
		return false,nil
	}
	d:=r.Data[0]
	if d.Code!=200{
		return false,nil
	}
	if len(d.Results)==0{
		return false,nil
	}
	for _,v:=range d.Results{
		if v.Suggestion=="block"{
			return false,nil
		}
	}
	return true,nil
}

func (this_* DefaultClient)CheckImage(dataID string,url string,log*logger.Logger)(success bool,err error)  {
	path := "/green/image/scan"
	
	clientInfo := &ClientInfo{}
	
	// request data
	bizType := "Green"
	scenes := []string{"porn","terrorism","ad","qrcode"}
	
	task := &Task{DataId: dataID, Url:url}
	tasks := []*Task{task}
	
	bizData := &BizData{bizType, scenes, tasks}
	resp:=this_.GetResponse(path, clientInfo, bizData,log)
	defer func() {
		if !success || err!=nil{
			if log!=nil{
				log.Error("aligreen failed",zap.Error(err),zap.ByteString("resp",resp))
			}
		}
	}()
	r:=&Result{}
	err =json.Unmarshal(resp,r)
	if err!=nil{
		return false,err
	}
	if r.Code!=200{
		return false,nil
	}
	if len(r.Data)!=1 {
		return false,nil
	}
	d:=r.Data[0]
	if d.Code!=200{
		return false,nil
	}
	if len(d.Results)==0{
		return false,nil
	}
	for _,v:=range d.Results{
		if v.Suggestion!="pass"{
			return false,nil
		}
	}
	return true,nil
}

func CheckText(dataID string,text string,log*logger.Logger)(success bool,err error)  {
	return defaultClient.CheckText(dataID,text,log)
}

func CheckImage(dataID string,url string,log*logger.Logger)(success bool,err error)  {
	return defaultClient.CheckImage(dataID,url,log)
}

type IAliYunClient interface {
	GetResponse(path string, clientInfo *ClientInfo, bizData *BizData) []byte
}
