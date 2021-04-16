package consulclient

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/njmdk/common/logger"
	"github.com/njmdk/common/utils"
)

func TestConsulClient(t *testing.T) {
	r := require.New(t)
	log, err := logger.New(utils.NewLogFileName("TestConsulClient", ""), "./log", zap.InfoLevel, false)
	r.NoError(err)
	cc, err := NewConsulClient("127.0.0.1:8500", "test", log)
	r.NoError(err)
	cfgAuth, err := cc.GetMysql("mysql.auth")
	r.NoError(err)
	r.Equal("root", cfgAuth.User)
	cfgAuth, err = cc.GetMysql("mysql.auth1")
	r.Error(err)
	value, err := cc.GetKey("kg_proxy_addr1")
	r.Equal(ErrorNotFound, err)
	fmt.Println(string(value))
	value, err = cc.GetKey("kg_proxy_addr")
	r.Equal(nil, err)
	fmt.Println(string(value))
}
