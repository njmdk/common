package mysql_uuid

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/njmdk/common/db"
	"github.com/njmdk/common/logger"
	"github.com/njmdk/common/utils"
)

func TestCreateUUIDCreator(t *testing.T) {
	r := require.New(t)
	log, err := logger.New("test_mysql_uuid", "./log", zap.InfoLevel, false)
	if err != nil {
		panic(err)
	}
	mysqlOpenAPI := db.NewMySQL(&db.Config{
		Addr:         "127.0.0.1:3306",
		User:         "root",
		Password:     "root",
		Database:     "ag",
		MaxOpenConnS: 10,
		MaxIdleConnS: 0,
		MaxLifeTime:  utils.Duration{Duration: time.Minute},
	}, log)
	uuid, err := CreateUUIDCreator(mysqlOpenAPI)
	r.NoError(err)
	for i := 0; i < 2; i++ {
		uid := uuid.Next()
		fmt.Println(uid)
		fmt.Println(uuid.NextString())
		fmt.Println(uuid.NextBase34())
		fmt.Println(uuid.NextBase62())
		fmt.Println(uuid.NextString())
		uu := uuid.Next16String()
		fmt.Println(uu)
		uu = uuid.Next32String()
		fmt.Println(uu)
	}
	//_, err = mysqlOpenAPI.Exec("INSERT INTO `ag`.`log_money` (`order_id`, `user_name`, `uuid`, `account`, `prefix`, `transfer_type`, `amount`, `created_time`, `status`, `complete_time`, `stat`) VALUES ('10010vJ6eI3wj6H', '10010eDIS2', 'skovsi3hvugwysm1', 'mofang', '1001', 'IN', '1.00', '2020-03-19 16:30:48', '1', '2020-03-19 16:30:48', '1');")
	//r.NoError(err)
}
