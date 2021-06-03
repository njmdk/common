package wx_login

import (
	"fmt"
	"github.com/njmdk/common/logger"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestWX_GetAccessToken(t *testing.T) {
	log, err := logger.InitDefaultLogger("test getaccesstoken", "log", zap.DebugLevel, false)
	require.NoError(t, err)
	client, err := NewWXLogin("","",log)
	require.NoError(t, err)
	result, err := client.GetAccessToken("","wx")
	require.NoError(t, err)
	fmt.Println(result)
}

func TestWX_GetUserInfo(t *testing.T) {
	log, err := logger.InitDefaultLogger("test getUserInfo", "log", zap.DebugLevel, false)
	require.NoError(t, err)
	client, err := NewWXLogin("","",log)
	require.NoError(t, err)
	result, err := client.GetUserInfo("", "")
	require.NoError(t, err)
	fmt.Println(result)
}
