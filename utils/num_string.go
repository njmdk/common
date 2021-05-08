package utils

import (
	"strconv"

	"github.com/njmdk/common/crypt"
)

func EncodeNumString(src string) string {
	high := src[:12]
	low := ""
	if len(src) > 12 {
		low = src[12:]
	}
	ih, err := strconv.Atoi(high)
	if err != nil {
		return src
	}
	il, err := strconv.Atoi(low)
	if err != nil {
		return src
	}
	return crypt.Base62Encode(int64(ih)) + "." + crypt.Base62Encode(int64(il))
}
