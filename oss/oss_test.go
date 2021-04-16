package oss

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var ossConfig = []byte(`{
  "access_key_id":"",
  "access_key_secret":"",
  "endpoint":"https://oss-cn-chengdu.aliyuncs.com",
  "bucket":""
}`)

var aliOss *OSS

func init() {
	o := &OSS{}
	err := json.Unmarshal(ossConfig, o)
	if err != nil {
		panic(err)
	}
	err = o.New()
	if err != nil {
		panic(err)
	}
	aliOss = o
}

func TestOSS_IsObjectExists(t *testing.T) {
	r := require.New(t)
	b, err := aliOss.IsObjectExists("league/1743.jpg")
	r.NoError(err)
	fmt.Println(b)
}

func TestOSS_CheckAndPut(t *testing.T) {
	r := require.New(t)
	err := aliOss.CheckAndPut("https://qn.feijing88.com/egame/csgo/league/3c458e59a1eaa4f28815fa0131c380f3.png", "league/1743.jpg")
	r.NoError(err)
}
