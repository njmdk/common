package utils

import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

func TestBase62(t *testing.T) {
	//r := require.New(t)
	//excepted := int64(216335982)
	//e := Base62EncodeMin6Max11(excepted)
	//fmt.Println(e)
	//actual := Base62Decode(e)
	//fmt.Println(excepted, actual)
	//r.Equal(excepted, actual)
	//actual = Base62DecodeMin6Max11(e)
	//fmt.Println(excepted, actual)
	//r.Equal(excepted, actual)
	num := Base62Decode("NJR0ebJ")
	a:=big.NewInt(math.MaxInt64)
	for i:=0;i<2;i++{
		a = a.Add(a,a)
	}
	fmt.Println(num,math.MaxInt64,a.String())
}
