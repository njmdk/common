package wx_pay

import (
    "bytes"
    "context"
    "crypto/hmac"
    "crypto/md5"
    "crypto/sha256"
    "crypto/tls"
    "encoding/base64"
    "encoding/hex"
    "encoding/pem"
    "encoding/xml"
    "errors"
    "fmt"
    "io/ioutil"
    "net"
    "net/http"
    "sort"
    "strconv"
    "strings"
    "time"
    
    "golang.org/x/crypto/pkcs12"
)

const (
    Fail                       = "FAIL"
    Success                    = "SUCCESS"
    HMACSHA256                 = "HMAC-SHA256"
    MD5                        = "MD5"
    Sign                       = "sign"
    MicroPayUrl                = "https://api.mch.weixin.qq.com/pay/micropay"
    UnifiedOrderUrl            = "https://api.mch.weixin.qq.com/pay/unifiedorder"
    OrderQueryUrl              = "https://api.mch.weixin.qq.com/pay/orderquery"
    ReverseUrl                 = "https://api.mch.weixin.qq.com/secapi/pay/reverse"
    CloseOrderUrl              = "https://api.mch.weixin.qq.com/pay/closeorder"
    RefundUrl                  = "https://api.mch.weixin.qq.com/secapi/pay/refund"
    RefundQueryUrl             = "https://api.mch.weixin.qq.com/pay/refundquery"
    DownloadBillUrl            = "https://api.mch.weixin.qq.com/pay/downloadbill"
    DownloadFundFlowUrl        = "https://api.mch.weixin.qq.com/pay/downloadfundflow"
    ReportUrl                  = "https://api.mch.weixin.qq.com/payitil/report"
    ShortUrl                   = "https://api.mch.weixin.qq.com/tools/shorturl"
    AuthCodeToOpenidUrl        = "https://api.mch.weixin.qq.com/tools/authcodetoopenid"
    SandboxMicroPayUrl         = "https://api.mch.weixin.qq.com/sandboxnew/pay/micropay"
    SandboxUnifiedOrderUrl     = "https://api.mch.weixin.qq.com/sandboxnew/pay/unifiedorder"
    SandboxOrderQueryUrl       = "https://api.mch.weixin.qq.com/sandboxnew/pay/orderquery"
    SandboxReverseUrl          = "https://api.mch.weixin.qq.com/sandboxnew/secapi/pay/reverse"
    SandboxCloseOrderUrl       = "https://api.mch.weixin.qq.com/sandboxnew/pay/closeorder"
    SandboxRefundUrl           = "https://api.mch.weixin.qq.com/sandboxnew/secapi/pay/refund"
    SandboxRefundQueryUrl      = "https://api.mch.weixin.qq.com/sandboxnew/pay/refundquery"
    SandboxDownloadBillUrl     = "https://api.mch.weixin.qq.com/sandboxnew/pay/downloadbill"
    SandboxDownloadFundFlowUrl = "https://api.mch.weixin.qq.com/sandboxnew/pay/downloadfundflow"
    SandboxReportUrl           = "https://api.mch.weixin.qq.com/sandboxnew/payitil/report"
    SandboxShortUrl            = "https://api.mch.weixin.qq.com/sandboxnew/tools/shorturl"
    SandboxAuthCodeToOpenidUrl = "https://api.mch.weixin.qq.com/sandboxnew/tools/authcodetoopenid"
)

const certData = ""

type Account struct {
    appID     string
    mchID     string
    apiKey    []byte
    certData  []byte
    isSandbox bool
}

// 创建微信支付账号
func NewWEChatAccount(appID string, mchID string, apiKey []byte,isSandbox bool) *Account {
    bd,_:=base64.StdEncoding.DecodeString(certData)
    return &Account{
        appID:  appID,
        mchID:  mchID,
        apiKey: apiKey,
        certData: bd,
        isSandbox: isSandbox,
    }
}

type Params map[string]string

// map本来已经是引用类型了，所以不需要 *Params
func (p Params) SetString(k, s string) Params {
    p[k] = s
    return p
}

func (p Params) GetString(k string) string {
    s, _ := p[k]
    return s
}

func (p Params) SetInt64(k string, i int64) Params {
    p[k] = strconv.FormatInt(i, 10)
    return p
}

func (p Params) GetInt64(k string) int64 {
    i, _ := strconv.ParseInt(p.GetString(k), 10, 64)
    return i
}

// 判断key是否存在
func (p Params) ContainsKey(key string) bool {
    _, ok := p[key]
    return ok
}

type Notifies struct{}

// 通知成功
func (n *Notifies) OK() string {
    var params = make(Params)
    params.SetString("return_code", Success)
    params.SetString("return_msg", "ok")
    return MapToXml(params)
}

// 通知不成功
func (n *Notifies) NotOK(errMsg string) string {
    var params = make(Params)
    params.SetString("return_code", Fail)
    params.SetString("return_msg", errMsg)
    return MapToXml(params)
}

func XmlToMap(xmlStr string) Params {
    params := make(Params)
    decoder := xml.NewDecoder(strings.NewReader(xmlStr))
    
    var (
        key   string
        value string
    )
    
    for t, err := decoder.Token(); err == nil; t, err = decoder.Token() {
        switch token := t.(type) {
        case xml.StartElement: // 开始标签
            key = token.Name.Local
        case xml.CharData: // 标签内容
            content := string([]byte(token))
            value = content
        }
        if key != "xml" {
            if value != "\n" {
                params.SetString(key, value)
            }
        }
    }
    
    return params
}

func MapToXml(params Params) string {
    var buf bytes.Buffer
    buf.WriteString(`<xml>`)
    for k, v := range params {
        buf.WriteString(`<`)
        buf.WriteString(k)
        buf.WriteString(`><![CDATA[`)
        buf.WriteString(v)
        buf.WriteString(`]]></`)
        buf.WriteString(k)
        buf.WriteString(`>`)
    }
    buf.WriteString(`</xml>`)
    
    return buf.String()
}

// 用时间戳生成随机字符串
func NonceStr() string {
    return strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
}

// 将Pkcs12转成Pem
func pkcs12ToPem(p12 []byte, password string) (c *tls.Certificate, e error) {
    blocks, err := pkcs12.ToPEM(p12, password)
    if err != nil {
        return nil, err
    }
    // 从恐慌恢复
    defer func() {
        if x := recover(); x != nil {
            c = nil
            e = fmt.Errorf("panic: %v",x)
        }
    }()
    
 
    
    var pemData []byte
    for _, b := range blocks {
        pemData = append(pemData, pem.EncodeToMemory(b)...)
    }
    
    cert, err := tls.X509KeyPair(pemData, pemData)
    if err != nil {
       return nil,err
    }
    return &cert,nil
}

const bodyType = "application/xml; charset=utf-8"

type WeChatClient struct {
    account              *Account // 支付账号
    signType             string   // 签名类型
    httpConnectTimeoutMs int      // 连接超时时间
    httpReadTimeoutMs    int      // 读取超时时间
    client * http.Client
}

// 创建微信支付客户端
func NewWEChatClient(account *Account) *WeChatClient {
    return &WeChatClient{
        account:              account,
        signType:             MD5,
        httpConnectTimeoutMs: 2000,
        httpReadTimeoutMs:    1000,
        client: &http.Client{
            Transport: &http.Transport{
                Proxy: nil,
                DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
                    return net.DialTimeout(network, addr, time.Second*10)
                },
                MaxIdleConns:        10,
                MaxIdleConnsPerHost: 1000,
                IdleConnTimeout:     time.Minute * 5,
                TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
            },
            Timeout: time.Second*10,
            CheckRedirect: func(req *http.Request, via []*http.Request) error {
                if len(via) >= 100 {
                    return errors.New("stopped after 100 redirects")
                }
                return nil
            },
        },
    }
}

func (c *WeChatClient) SetHttpConnectTimeoutMs(ms int) {
    c.httpConnectTimeoutMs = ms
}

func (c *WeChatClient) SetHttpReadTimeoutMs(ms int) {
    c.httpReadTimeoutMs = ms
}

func (c *WeChatClient) SetSignType(signType string) {
    c.signType = signType
}

func (c *WeChatClient) SetAccount(account *Account) {
    c.account = account
}

// 向 params 中添加 appid、mch_id、nonce_str、sign_type、sign
func (c *WeChatClient) fillRequestData(params Params) Params {
    params["appid"] = c.account.appID
    params["mch_id"] = c.account.mchID
    params["nonce_str"] = NonceStr()
    params["sign_type"] = c.signType
    params["sign"] = c.Sign(params)
    return params
}
// 向 params 中添加 appid、mch_id、nonce_str、sign_type、sign
func (c *WeChatClient) FillJSAPIRequestData(params Params) Params {
    params["appId"] = c.account.appID
    params["timeStamp"] = fmt.Sprintf("%d",time.Now().Unix())
    params["signType"] = c.signType
    params["paySign"] = c.SignJSAPI(params)
    return params
}

// https no cert post
func (c *WeChatClient) postWithoutCert(url string, params Params) (string, error) {
    p := c.fillRequestData(params)
    response, err := c.client.Post(url, bodyType, strings.NewReader(MapToXml(p)))
    if err != nil {
        return "", err
    }
    if response.StatusCode!=http.StatusOK{
        return "",fmt.Errorf("%d,%s",response.StatusCode,response.Status)
    }
    defer response.Body.Close()
    res, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return "", err
    }
    return string(res), nil
}

// https need cert post
func (c *WeChatClient) postWithCert(url string, params Params) (string, error) {
    if c.account.certData == nil {
        return "", errors.New("证书数据为空")
    }
    
    // 将pkcs12证书转成pem
    cert,err := pkcs12ToPem(c.account.certData, c.account.mchID)
    if err!=nil{
        return "", err
    }
    
    config := &tls.Config{
        Certificates: []tls.Certificate{*cert},
    }
    transport := &http.Transport{
        TLSClientConfig:    config,
        DisableCompression: true,
    }
    h := &http.Client{Transport: transport}
    p := c.fillRequestData(params)
    response, err := h.Post(url, bodyType, strings.NewReader(MapToXml(p)))
    if err != nil {
        return "", err
    }
    defer response.Body.Close()
    res, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return "", err
    }
    return string(res), nil
}

// 生成带有签名的xml字符串
func (c *WeChatClient) generateSignedXml(params Params) string {
    sign := c.Sign(params)
    params.SetString(Sign, sign)
    return MapToXml(params)
}

// 验证签名
func (c *WeChatClient) ValidSign(params Params) bool {
    if !params.ContainsKey(Sign) {
        return false
    }
    return params.GetString(Sign) == c.Sign(params)
}

// 签名
func (c *WeChatClient) Sign(params Params) string {
    // 创建切片
    var keys = make([]string, 0, len(params))
    // 遍历签名参数
    for k := range params {
        if k != "sign" { // 排除sign字段
            keys = append(keys, k)
        }
    }
    
    // 由于切片的元素顺序是不固定，所以这里强制给切片元素加个顺序
    sort.Strings(keys)
    
    //创建字符缓冲
    var buf bytes.Buffer
    for _, k := range keys {
        if len(params.GetString(k)) > 0 {
            buf.WriteString(k)
            buf.WriteString(`=`)
            buf.WriteString(params.GetString(k))
            buf.WriteString(`&`)
        }
    }
    // 加入apiKey作加密密钥
    buf.WriteString(`key=`)
    buf.WriteString(string(c.account.apiKey))
    
    var (
        dataMd5    [16]byte
        dataSha256 []byte
        str        string
    )
    
    switch c.signType {
    case MD5:
        dataMd5 = md5.Sum(buf.Bytes())
        str = hex.EncodeToString(dataMd5[:]) //需转换成切片
    case HMACSHA256:
        h := hmac.New(sha256.New, []byte(c.account.apiKey))
        h.Write(buf.Bytes())
        dataSha256 = h.Sum(nil)
        str = hex.EncodeToString(dataSha256[:])
    }
    
    return strings.ToUpper(str)
}

// 签名
func (c *WeChatClient) SignJSAPI(params Params) string {
    // 创建切片
    var keys = make([]string, 0, len(params))
    // 遍历签名参数
    for k := range params {
        if k != "paySign" { // 排除sign字段
            keys = append(keys, k)
        }
    }

    // 由于切片的元素顺序是不固定，所以这里强制给切片元素加个顺序
    sort.Strings(keys)

    //创建字符缓冲
    var buf bytes.Buffer
    for _, k := range keys {
        if len(params.GetString(k)) > 0 {
            buf.WriteString(k)
            buf.WriteString(`=`)
            buf.WriteString(params.GetString(k))
            buf.WriteString(`&`)
        }
    }
    // 加入apiKey作加密密钥
    buf.WriteString(`key=`)
    buf.WriteString(string(c.account.apiKey))

    var (
        dataMd5    [16]byte
        dataSha256 []byte
        str        string
    )

    switch c.signType {
    case MD5:
        dataMd5 = md5.Sum(buf.Bytes())
        str = hex.EncodeToString(dataMd5[:]) //需转换成切片
    case HMACSHA256:
        h := hmac.New(sha256.New, []byte(c.account.apiKey))
        h.Write(buf.Bytes())
        dataSha256 = h.Sum(nil)
        str = hex.EncodeToString(dataSha256[:])
    }

    return strings.ToUpper(str)
}
// 处理 HTTPS API返回数据，转换成Map对象。return_code为SUCCESS时，验证签名。
func (c *WeChatClient) processResponseXml(xmlStr string) (Params, error) {
    var returnCode string
    params := XmlToMap(xmlStr)
    if params.ContainsKey("return_code") {
        returnCode = params.GetString("return_code")
    } else {
        return nil, errors.New("no return_code in XML")
    }
    if returnCode == Fail {
        return params, nil
    } else if returnCode == Success {
        if c.ValidSign(params) {
            return params, nil
        } else {
            return nil, errors.New("invalid sign value in XML")
        }
    } else {
        return nil, errors.New("return_code value is invalid in XML")
    }
}

// 统一下单
func (c *WeChatClient) UnifiedOrder(params Params) (Params, error) {
    var url string
    if c.account.isSandbox {
        url = SandboxUnifiedOrderUrl
    } else {
        url = UnifiedOrderUrl
    }
    xmlStr, err := c.postWithoutCert(url, params)
    if err != nil {
        return nil, err
    }
    return c.processResponseXml(xmlStr)
}

// 刷卡支付
func (c *WeChatClient) MicroPay(params Params) (Params, error) {
    var url string
    if c.account.isSandbox {
        url = SandboxMicroPayUrl
    } else {
        url = MicroPayUrl
    }
    xmlStr, err := c.postWithoutCert(url, params)
    if err != nil {
        return nil, err
    }
    return c.processResponseXml(xmlStr)
}

// 退款
func (c *WeChatClient) Refund(params Params) (Params, error) {
    var url string
    if c.account.isSandbox {
        url = SandboxRefundUrl
    } else {
        url = RefundUrl
    }
    xmlStr, err := c.postWithCert(url, params)
    if err != nil {
        return nil, err
    }
    return c.processResponseXml(xmlStr)
}

// 订单查询
func (c *WeChatClient) OrderQuery(params Params) (Params, error) {
    var url string
    if c.account.isSandbox {
        url = SandboxOrderQueryUrl
    } else {
        url = OrderQueryUrl
    }
    xmlStr, err := c.postWithoutCert(url, params)
    if err != nil {
        return nil, err
    }
    return c.processResponseXml(xmlStr)
}

// 退款查询
func (c *WeChatClient) RefundQuery(params Params) (Params, error) {
    var url string
    if c.account.isSandbox {
        url = SandboxRefundQueryUrl
    } else {
        url = RefundQueryUrl
    }
    xmlStr, err := c.postWithoutCert(url, params)
    if err != nil {
        return nil, err
    }
    return c.processResponseXml(xmlStr)
}

// 撤销订单
func (c *WeChatClient) Reverse(params Params) (Params, error) {
    var url string
    if c.account.isSandbox {
        url = SandboxReverseUrl
    } else {
        url = ReverseUrl
    }
    xmlStr, err := c.postWithCert(url, params)
    if err != nil {
        return nil, err
    }
    return c.processResponseXml(xmlStr)
}

// 关闭订单
func (c *WeChatClient) CloseOrder(params Params) (Params, error) {
    var url string
    if c.account.isSandbox {
        url = SandboxCloseOrderUrl
    } else {
        url = CloseOrderUrl
    }
    xmlStr, err := c.postWithoutCert(url, params)
    if err != nil {
        return nil, err
    }
    return c.processResponseXml(xmlStr)
}

// 对账单下载
func (c *WeChatClient) DownloadBill(params Params) (Params, error) {
    var url string
    if c.account.isSandbox {
        url = SandboxDownloadBillUrl
    } else {
        url = DownloadBillUrl
    }
    xmlStr, err := c.postWithoutCert(url, params)
    
    p := make(Params)
    
    // 如果出现错误，返回XML数据
    if strings.Index(xmlStr, "<") == 0 {
        p = XmlToMap(xmlStr)
        return p, err
    } else { // 正常返回csv数据
        p.SetString("return_code", Success)
        p.SetString("return_msg", "ok")
        p.SetString("data", xmlStr)
        return p, err
    }
}

func (c *WeChatClient) DownloadFundFlow(params Params) (Params, error) {
    var url string
    if c.account.isSandbox {
        url = SandboxDownloadFundFlowUrl
    } else {
        url = DownloadFundFlowUrl
    }
    xmlStr, err := c.postWithCert(url, params)
    
    p := make(Params)
    
    // 如果出现错误，返回XML数据
    if strings.Index(xmlStr, "<") == 0 {
        p = XmlToMap(xmlStr)
        return p, err
    } else { // 正常返回csv数据
        p.SetString("return_code", Success)
        p.SetString("return_msg", "ok")
        p.SetString("data", xmlStr)
        return p, err
    }
}

// 交易保障
func (c *WeChatClient) Report(params Params) (Params, error) {
    var url string
    if c.account.isSandbox {
        url = SandboxReportUrl
    } else {
        url = ReportUrl
    }
    xmlStr, err := c.postWithoutCert(url, params)
    if err != nil {
        return nil, err
    }
    return c.processResponseXml(xmlStr)
}

// 转换短链接
func (c *WeChatClient) ShortUrl(params Params) (Params, error) {
    var url string
    if c.account.isSandbox {
        url = SandboxShortUrl
    } else {
        url = ShortUrl
    }
    xmlStr, err := c.postWithoutCert(url, params)
    if err != nil {
        return nil, err
    }
    return c.processResponseXml(xmlStr)
}

// 授权码查询OPENID接口
func (c *WeChatClient) AuthCodeToOpenid(params Params) (Params, error) {
    var url string
    if c.account.isSandbox {
        url = SandboxAuthCodeToOpenidUrl
    } else {
        url = AuthCodeToOpenidUrl
    }
    xmlStr, err := c.postWithoutCert(url, params)
    if err != nil {
        return nil, err
    }
    return c.processResponseXml(xmlStr)
}