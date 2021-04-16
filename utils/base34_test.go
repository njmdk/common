package utils

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common/math"

	"github.com/njmdk/common/timer"
)

func TestBase34(t *testing.T) {
	id := (math.MaxInt64) << 16 >> 16
	c := Base34(uint64(id))
	fmt.Println(c)
	fmt.Println(strconv.Itoa(id))

	now := timer.Now()
	now = now.AddDate(0, 1, 0)
	fmt.Println(timer.ToString(now))
}
