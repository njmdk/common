package http

import (
	"encoding/xml"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/njmdk/common/logger"
)

type XmlT struct {
	XMLName xml.Name `xml:"result"`
	Info    string   `xml:"info,attr"`
	Msg     string   `xml:"msg,attr"`
}

type Result1 struct {
}

func TestHttpClient_Do(t *testing.T) {
	r := require.New(t)
	log, _ := logger.InitDefaultLogger("TestHttpClient_Do", ".", zap.DebugLevel, true)

	client, err := NewHTTPClientWithTimeout("https://www.douyu.com", nil, log, false, time.Second*5)
	r.NoError(err)
	//queryValues := utils.NewURLValues().Add("params", "MjgzMzY0NzI1NHxhYWFhcGlhZ3xDfA").Add("sign", "eff5fa002770eb3bfe0d7654271e4b4f")
	//out := map[string]interface{}{}
	data, err := client.GET("/522425",
		nil,
		//utils.NewHTTPHeader().Add("User-Agent", "WEB_LIB_GI_GR4_AGIN"),
		nil,
		nil)
	r.NoError(err)
	fmt.Println(string(data))
}
