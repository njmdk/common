package workpool

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
	
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/njmdk/common/logger"
)

func TestNewWorkPool(t *testing.T) {
	r := require.New(t)
	log, _ := logger.InitDefaultLogger("TestNewWorkPool", ".", zap.InfoLevel, false)
	testCount := int64(2000000)
	workPool := NewWorkPool(0, log)
	workPool.Run(nil)
	var count int64
	n:=time.Now()
	for i := int64(0); i < testCount; i++ {
		workPool.Post(func() {
			atomic.AddInt64(&count, 1)
		})
	}
	workPool.Close()
	fmt.Println(time.Now().Sub(n))

	workPool.WaitClosed()
	r.Equal(testCount, count)
}
