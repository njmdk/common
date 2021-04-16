package aligreen

import (
    "fmt"
    "testing"
    "time"
    
    "github.com/google/uuid"
)

func TestDefaultClient_GetResponse(t *testing.T) {
    s:=time.Now()
    fmt.Println(CheckText(uuid.New().String(),"123",nil))
    fmt.Println(time.Now().Sub(s))
}

func TestDefaultClient_GetResponse1(t *testing.T) {
    s:=time.Now()
    fmt.Println(defaultClient.CheckImage(uuid.New().String(),"https://djtest2.oss-cn-chengdu.aliyuncs.com/chat/2e0b7d18-b13a-46e9-91d7-e3ccfff5f360.jpeg",nil))
    fmt.Println(time.Now().Sub(s))
}