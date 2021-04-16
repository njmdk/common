package logger

import (
	"fmt"
	"testing"
)

func TestDeleteBeforeNDay(t *testing.T) {
	fmt.Println(reg.FindAllString("2020_1_22312", -1))
	fmt.Println(reg.FindAllString("2020_01_21231", -1))
	fmt.Println(reg.FindAllString("2020_1_0212312", -1))
	fmt.Println(reg.FindAllString("2020_1_02123123", -1))
}
