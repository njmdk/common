package ali_pay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"hash"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AliPay struct {
	AppId string `json:"app_id"`
	AppPriKey []byte `json:"app_pri_key"`
	AppPubKey []byte `json:"app_pub_key"`
	CallBack string `json:"call_back"`
	QuitUrl string `json:"quit_url"`
}
func NewAliPay(appId string,appPriKey []byte,appPubKey []byte,callback string,quitUrl string) (*AliPay,error) {
	aliPay := &AliPay{
		AppId:  appId,
		AppPriKey: appPriKey,
		AppPubKey: appPubKey,
		CallBack:callback,
		QuitUrl:quitUrl,
	}
	return aliPay,nil
}

func (this_ *AliPay) GenBizContent(subject, body, outTradeNo, productCode string, totalAmount float64,quitUrl string) (string, error) {
	m := make(map[string]interface{})
	m["out_trade_no"] = outTradeNo
	m["product_code"] = productCode
	m["total_amount"] = totalAmount //TODO
	m["subject"] = subject
	m["body"] = body
	m["timeout_express"] = "30m"
	if quitUrl != "" {
		m["quit_url"] = quitUrl
	}

	jsonStr, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(jsonStr), nil
}
func (this_ *AliPay) FillSign2Data(bizContent string, privateKey []byte, method string, signType string,returnUrl string) (url.Values, error) {
	data := url.Values{}
	data.Set("app_id", this_.AppId)
	data.Set("method", method)
	data.Set("charset", "utf-8")
	data.Set("sign_type", signType)
	now := time.Now().Format("2006-01-02 15:04:05")
	if method == "alipay.trade.wap.pay" {
		data.Set("format", "json")
	}
	if returnUrl != "" {
		data.Set("return_url", returnUrl)
	}
	data.Set("timestamp", now)
	data.Set("version", "1.0")
	data.Set("notify_url", this_.CallBack + PathCallback)
	data.Set("biz_content", bizContent)

	//生成签名
	signContentBytes, _ := url.QueryUnescape(data.Encode())
	//fmt.Println("data to be signed", signContentBytes)
	signature, err := this_.Sign([]byte(signContentBytes), signType, privateKey)
	if err != nil {
		return nil, err
	}
	data.Set("sign", signature)
	return data, nil
}
func (this_ *AliPay)Sign(data []byte, SignType string, pemPriKey []byte) (signature string, err error) {
	var h hash.Hash
	var hType crypto.Hash
	switch SignType {
	case RSA:
		h = sha1.New()
		hType = crypto.SHA1
	case RSA2:
		h = sha256.New()
		hType = crypto.SHA256
	}
	h.Write(data)
	d := h.Sum(nil)
	pk, err := this_.ParsePrivateKey(pemPriKey)
	if err != nil {
		return "",err
	}
	bs, err := rsa.SignPKCS1v15(rand.Reader, pk, hType, d)
	//bs, err := rsa.SignPSS(rand.Reader, pk, hType, d,nil)

	if err != nil {
		return "",err
	}
	signature = base64.StdEncoding.EncodeToString(bs)
	return signature, nil
}

func (this_ *AliPay)ParsePrivateKey(privateKey []byte) (pk *rsa.PrivateKey, err error) {
	block, _ := pem.Decode(privateKey)
	//key, err := base64.StdEncoding.DecodeString(privateKey)
	//if err != nil {
	//	return nil,err
	//}
	//pk, err := x509.ParsePKCS1PrivateKey(key)
	//if err != nil {
	//	return nil,err
	//}
	if block == nil {
		return nil,errors.New("私钥格式错误1")
	}
	switch block.Type {
	case "RSA PRIVATE KEY":
		rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil,err
		}
		pk = rsaPrivateKey
		return pk,nil
	default:
		err = errors.New(fmt.Sprintf("私钥格式错误:%s", privateKey))
		return nil,err
	}
}
const (
	PEMBEGIN = "-----BEGIN RSA PRIVATE KEY-----\n"
	PEMEND = "\n-----END RSA PRIVATE KEY-----"
	PUBPEMBEGIN = "-----BEGIN PUBLIC KEY-----\n"
	PUBPEMEND = "\n-----END PUBLIC KEY-----"
)
// Rsa2PubSign RSA2公钥验证签名
func (this_ *AliPay)Rsa2PubSign(signContent, sign, publicKey string, hash crypto.Hash) bool {
	hashed := sha256.Sum256([]byte(signContent))
	pubKey, err := this_.ParsePublicKey(publicKey)
	if err != nil {
		return false
	}
	sig, _ := base64.StdEncoding.DecodeString(sign)
	err = rsa.VerifyPKCS1v15(pubKey, hash, hashed[:], sig)
	if err != nil {
		return false
	}
	return true
}

func (this_ *AliPay)ParsePublicKey(publicKey string) (*rsa.PublicKey, error) {
	publicKey = this_.FormatPublicKey(publicKey)
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, errors.New("公钥信息错误！")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pubKey.(*rsa.PublicKey), nil
}
// FormatPublicKey 组装公钥
func (this_ *AliPay)FormatPublicKey(publicKey string) string {
	if !strings.HasPrefix(publicKey, PUBPEMBEGIN) {
		publicKey = PUBPEMBEGIN + publicKey
	}
	if !strings.HasSuffix(publicKey, PUBPEMEND) {
		publicKey = publicKey + PUBPEMEND
	}
	return publicKey
}
func (this_ *AliPay) PostData(in url.Values) ([]byte,error) {
	data,_ := url.QueryUnescape(in.Encode())
	tmp := host + "?" +data
	fmt.Println(tmp)
	//resp,err := http.PostForm(host,in)
	resp,err := http.Post(host,"application/x-www-form-urlencoded;charset=utf-8",strings.NewReader(in.Encode()))
	if err != nil {
		return nil,err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body,err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil,err
		}
		return body,nil
	} else {
		return nil,errors.New(fmt.Sprintf("错误码 %d Status %s",resp.StatusCode,resp.Status))
	}
}
func (this_ *AliPay) GetData(in url.Values) ([]byte,error) {
	//data,err := json.Marshal(in)
	//if err != nil {
	//	return nil,err
	//}
	resp,err := http.Get(host + "?" + in.Encode())
	if err != nil {
		return nil,err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body,err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil,err
		}
		return body,nil
	} else {
		return nil,errors.New(fmt.Sprintf("错误码 %d Status %s",resp.StatusCode,resp.Status))
	}
}
func (this_ *AliPay) GenerateOrder(subject string,body string,orderId string,money float64,returnUrl string) (url.Values, error) {
	bizContent,err := this_.GenBizContent(subject,body,orderId,ProductFastTradePay,money,"")
	if err != nil {
		return nil,err
	}
	signData,err := this_.FillSign2Data(bizContent,this_.AppPriKey,alipayTradePagePay,RSA2,returnUrl)
	if err != nil {
		return nil,err
	}
	//resp,err := this_.PostData(signData)
	//if err != nil {
	//	return "",err
	//}
	//fmt.Println(resp)
	return signData,nil
}
type AliPayResponse struct {
	AliPayTradePreCreateResponse *PreCreateResponse `json:"alipay_trade_precreate_response"`
}
type PreCreateResponse struct {
	Code string `json:"code"`
	Msg string `json:"msg"`
	OutTradeNo string `json:"out_trade_no"`
	QrCode string `json:"qr_code"`
	Sign string `json:"sign"`
}
func (this_ *AliPay) PreCreateOrder(subject string,body string,orderId string,money float64,returnUrl string) (*AliPayResponse, error) {
	bizContent,err := this_.GenBizContent(subject,body,orderId,ProductFaceToFace,money,"")
	if err != nil {
		return nil,err
	}
	signData,err := this_.FillSign2Data(bizContent,this_.AppPriKey,alipayTradePreCreate,RSA2,returnUrl)
	if err != nil {
		return nil,err
	}
	resp,err := this_.PostData(signData)
	if err != nil {
		return nil,err
	}
	result := &AliPayResponse{}
	err = json.Unmarshal(resp,result)
	if err != nil {
		return nil,err
	}
	if result.AliPayTradePreCreateResponse.Code == "10000" {
		return result,nil
	} else {
		return nil,errors.New(fmt.Sprintf("错误码:%s 信息:%s",result.AliPayTradePreCreateResponse.Code,result.AliPayTradePreCreateResponse.Msg))
	}
}
func (this_ *AliPay) GenerateWapOrder(subject string,body string,orderId string,money float64,returnUrl string) (url.Values, error) {
	bizContent,err := this_.GenBizContent(subject,body,orderId,ProductQuickWapWay,money,this_.QuitUrl)
	if err != nil {
		return nil,err
	}
	signData,err := this_.FillSign2Data(bizContent,this_.AppPriKey,alipayTradeWapPay,RSA2,returnUrl)
	if err != nil {
		return nil,err
	}
	//return signData,nil
	//resp,err := this_.PostData(signData)
	//if err != nil {
	//	return nil,err
	//}
	//fmt.Println(resp)
	return signData,nil
	//result := &AliPayResponse{}
	//err = json.Unmarshal(resp,result)
	//if err != nil {
	//	this_.log.Error("json.Unmarshal err",zap.Error(err))
	//	return nil,err
	//}
	//if result.AliPayTradePreCreateResponse.Code == "10000" {
	//	return result,nil
	//} else {
	//	return nil,errors.New(fmt.Sprintf("错误码:%s 信息:%s",result.AliPayTradePreCreateResponse.Code,result.AliPayTradePreCreateResponse.Msg))
	//}
}
func (this_ *AliPay) ProcessNotify(postForm url.Values) error {
	var sign string
	var signType string
	params := url.Values{}
	mapParams := make(map[string]string)
	for k,v := range postForm {
		if len(v) > 0 {
			mapParams[k] = v[0]
			if k == "sign" {
				sign = v[0]
			} else if k == "sign_type" {
				signType = v[0]
			} else {
				params[k] = v
			}
		}
	}
	signText,err := url.QueryUnescape(params.Encode())
	if err != nil {
		return err
	}
	if signType == RSA2 {
		ret := this_.Rsa2PubSign(signText,sign,string(this_.AppPubKey),crypto.SHA256)
		if !ret {
			return errors.New("Rsa2PubSign错误")
		}
	} else {
		return errors.New("不支持的加密")
	}

	tradeStatus,ok := mapParams["trade_status"]
	if !ok {
		return errors.New("不存在的字段 trade_status")
	}
	if tradeStatus == "TRADE_SUCCESS" {
		appId,ok := mapParams["app_id"]
		if !ok {
			return errors.New("不存在的字段 app_id")
		}
		if appId != this_.AppId {
			return errors.New("app_id不匹配")
		}
		return nil
	} else {
		return errors.New(fmt.Sprintf("交易状态错误 %s",tradeStatus))
	}
}