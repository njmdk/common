package ali_pay

const (
	ProductFastTradePay = "FAST_INSTANT_TRADE_PAY"
	ProductFaceToFace = "FACE_TO_FACE_PAYMENT"
	ProductQuickWapWay = "QUICK_WAP_WAY"

	alipayTradePagePay = "alipay.trade.page.pay"
	alipayTradePreCreate = "alipay.trade.precreate"
	alipayTradeWapPay = "alipay.trade.wap.pay"

	RSA = "RSA"
	RSA2 = "RSA2"
	host = "https://openapi.alipay.com/gateway.do"
	PathCallback = "/ali/callback"
	QuitUrl = "http://127.0.0.1"

	tradeSuccess = "TRADE_SUCCESS"
	failMsg = "fail"
	successMsg = "success"
)

type BizContent struct {
	OutTradeNo string `json:"out_trade_no"`
	ProductCode string `json:"product_code"`
	TotalAmount float64 `json:"total_amount"`
	Subject string `json:"subject"`
	Body string `json:"body"`
	TimeoutExpress string `json:"timeout_express"`
}
type SignDataInfo struct {
	AppId string `json:"app_id"`
	Method string `json:"method"`
	Charset string `json:"charset"`
	SignType string `json:"sign_type"`
	Sign string `json:"sign"`
	Timestamp string `json:"timestamp"`
	Version string `json:"version"`
	NotifyUrl string `json:"notify_url"`
	BizContent string `json:"biz_content"`
}