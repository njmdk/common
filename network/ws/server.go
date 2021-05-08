package ws

import (
	"sync/atomic"

	"github.com/njmdk/common/crypt"
)

type WebSocket struct {
	index int64
}

func (this_ *WebSocket) GenID() string {
	newIndex := atomic.AddInt64(&this_.index, 1)
	return crypt.Base34(uint64(newIndex))
}
