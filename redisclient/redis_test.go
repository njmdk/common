package redisclient

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/njmdk/common/utils"
)

func TestRedisClient(t *testing.T) {
	r := require.New(t)
	redisClient, err := NewRedisClientWithConfig(&Config{
		Addr:       "127.0.0.1:6379",
		Password:   "",
		DataBase:   2,
		Timeout:    utils.Duration{Duration: time.Second},
		MaxRetries: 2,
	})
	r.NoError(err)
	//err = redisClient.HDel("123", "321").Err()
	//r.Equal(redis.Nil, err)
	//cmdS, err := redisClient.Pipelined(func(pp redis.Pipeliner) error {
	//	pp.HGetAll("socketio_addr_2833647254")
	//	pp.HGetAll("socketio_addr_28336472514")
	//	pp.HGetAll("socketio_addr_4116080134")
	//	return nil
	//})
	//
	//r.NoError(err)
	//fmt.Printf("%+v \n", cmdS)
	//fmt.Printf("%+v \n", cmdS[0].Args())
	//for _, v := range cmdS {
	//	fmt.Printf("%+v \n", v.(*redis.StringStringMapCmd).Val())
	//}
	u, err := NewRedisUUID(redisClient)
	r.NoError(err)

	m := &sync.Map{}
	wg := &sync.WaitGroup{}
	for x := 0; x < 10; x++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				uuid := u.Next28()
				fmt.Println(uuid)
				_, found := m.LoadOrStore(uuid, struct {
				}{})
				if found {
					panic(uuid)
				}
			}
		}()
	}
	wg.Wait()
	//fmt.Println(utils.Base62Encode(math.MaxInt64))
}

type AA struct {
	A string `json:"a"`
}
