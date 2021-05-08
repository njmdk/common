package crypt

import (
	"crypto/md5"
	"fmt"
	"io"
)

func MD5String(str string) string {
	w := md5.New()
	_, _ = io.WriteString(w, str)

	return fmt.Sprintf("%x", w.Sum(nil))
}

func MD5StringUpper(str string) string {
	w := md5.New()
	_, _ = io.WriteString(w, str)

	return fmt.Sprintf("%X", w.Sum(nil))
}

func MD5Bytes(bs []byte) string {
	md5 := md5.Sum(bs)
	md5Text := fmt.Sprintf("%x", md5)
	return md5Text
}
