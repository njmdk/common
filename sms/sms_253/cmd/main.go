package main

import (
	"fmt"

	"github.com/njmdk/common/sms/sms_253"
)

func main() {
	s := sms_253.NewStruct253SMS()
	s.Info.Phone = "11111111111"
	s.Info.Msg = "【测试】您好！您本次的验证码是: 1234 "
	out, err := sms_253.Post253SMS(s)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#v\n", out)
}
