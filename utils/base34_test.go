package utils

import (
	"fmt"
	"math"
	"strconv"
	"testing"
	"time"
)

func TestBase34(t *testing.T) {
	id := (math.MaxInt64) << 16 >> 16
	c := Base34(uint64(id))
	fmt.Println(c)
	fmt.Println(strconv.Itoa(id))

	now := time.Now()
	now = now.AddDate(0, 1, 0)
	fmt.Println(now.Format("2006-01-02 15:04:05"))
}
