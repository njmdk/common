package wx_pay

import (
    "fmt"
    "testing"
    
    "github.com/stretchr/testify/require"
)

var weChatKey = []byte(``)
var  weChatCert = []byte(``)

func TestNewWEChatAccountClient(t *testing.T) {
   r:=require.New(t)
    account:= NewWEChatAccount("","",[]byte(""),false)
    client:=NewWEChatClient(account)
    params := make(Params)
    params.SetString("body", "test").
        SetString("out_trade_no", "test436577857").
        SetInt64("total_fee", 1).
        SetString("spbill_create_ip", "127.0.0.1").
        SetString("notify_url", "http://notify.objcoding.com/notify").
        SetString("trade_type", "NATIVE").
        SetString("product_id","1")
    p, err := client.UnifiedOrder(params)
    r.NoError(err)
    fmt.Println(p)
}
