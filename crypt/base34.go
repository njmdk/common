package crypt

import (
	"container/list"
)

var baseStr = "0123456789ABCDEFGHJKLMNPQRSTUVWXYZ"

func Base34(id uint64) string {
	quotient := id
	mod := uint64(0)
	l := list.New()
	for quotient != 0 {
		mod = quotient % 34
		quotient = quotient / 34
		l.PushFront(baseStr[int(mod)])
	}
	listLen := l.Len()
	if listLen >= 6 {
		res := make([]byte, 0, listLen)
		for i := l.Front(); i != nil; i = i.Next() {
			res = append(res, i.Value.(byte))
		}
		return string(res)
	}
	res := make([]byte, 0, 6)
	for i := 0; i < 6; i++ {
		if i < 6-listLen {
			res = append(res, baseStr[0])
		} else {
			res = append(res, l.Front().Value.(byte))
			l.Remove(l.Front())
		}

	}
	return string(res)
}
