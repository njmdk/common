package sms_253

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPost253SMS(t *testing.T) {
	r := require.New(t)
	s := NewStruct253SMS()
	s.Info.Phone = "11111111111"
	s.Info.Msg = "【测试】您好！您本次的验证码是: 1234 "
	out, err := Post253SMS(s)
	r.NoError(err)
	fmt.Printf("%#v\n", out)
}
