package ws

import (
	"encoding/json"
	"errors"
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/njmdk/common/logger"
	"github.com/njmdk/common/utils"
)

type packet struct {
	Index string
	Msg   proto.Message
}

type Session struct {
	id          string
	cond        *sync.Cond
	sendBuf     []*packet
	ws          *websocket.Conn
	once        *sync.Once
	onClose     func(c *Session, err error)
	onMessage   func(c *Session, p *packet)
	close       int32
	contextType string
	log         *logger.Logger
}

func NewClient(id string, ws *websocket.Conn, onClose func(c *Session, err error), onMessage func(c *Session, p *packet)) *Session {
	c := &Session{
		id:          id,
		cond:        sync.NewCond(&sync.Mutex{}),
		ws:          ws,
		once:        &sync.Once{},
		onClose:     onClose,
		onMessage:   onMessage,
		contextType: "application/json",
	}
	return c
}

func (this_ *Session) Close(err error) {
	this_.once.Do(func() {
		if atomic.AddInt32(&this_.close, 1) == 1 {
			_ = this_.ws.Close()
			if this_.onClose != nil {
				this_.onClose(this_, err)
			}
		}
	})
}

func (this_ *Session) Send() {

}

func (this_ *Session) send(index string, msg proto.Message) {
	if this_.isClosed() {
		return
	}
	this_.cond.L.Lock()
	this_.sendBuf = append(this_.sendBuf, &packet{Index: index, Msg: msg})
	this_.cond.L.Unlock()
	this_.cond.Signal()
}

func (this_ *Session) isClosed() bool {
	return atomic.LoadInt32(&this_.close) != 0
}

func (this_ *Session) doSend() {
	utils.SafeGO(func(i interface{}) {
		this_.log.Error("doSend panic", zap.Any("panic info", i))
		this_.Close(errors.New("panic"))
	}, func() {
		var err error
		defer func() {
			this_.Close(err)
		}()
		for !this_.isClosed() {
			this_.cond.L.Lock()
			for len(this_.sendBuf) == 0 {
				this_.cond.Wait()
				if this_.isClosed() {
					return
				}
			}

			writeBuf := this_.sendBuf
			this_.sendBuf = nil
			this_.cond.L.Unlock()

			for _, v := range writeBuf {
				if err = this_.writeOne(v); err != nil {
					return
				}
			}
		}
	})

}

func (this_ *Session) writeOne(p *packet) error {
	msgName := proto.MessageName(p.Msg)
	if msgName == "" {
		return nil
	}
	var data []byte
	if this_.contextType == "application/json" {
		data, _ = json.Marshal(p.Msg)
	} else {
		data, _ = proto.Marshal(p.Msg)
	}
	return this_.ws.WriteJSON([]string{p.Index, msgName, utils.BytesToString(data)})
}

func (this_ *Session) doRead() {
	utils.SafeGO(func(i interface{}) {
		this_.log.Error("doSend panic", zap.Any("panic info", i))
		this_.Close(errors.New("panic"))
	}, func() {
		var err error
		defer func() {
			this_.Close(err)
		}()
		for !this_.isClosed() {
			var data []string
			err = this_.ws.ReadJSON(&data)
			if err != nil {
				return
			}
			if len(data) != 3 {
				err = errors.New("invalid msg")
				return
			}
			t := proto.MessageType(data[1])
			if t == nil {
				err = errors.New("invalid msg")
				return
			}
			msg := reflect.New(t.Elem()).Interface().(proto.Message)
			if this_.contextType == "application/json" {
				err = json.Unmarshal(utils.StringToBytes(data[2]), msg)
			} else {
				err = proto.Unmarshal(utils.StringToBytes(data[2]), msg)
			}
			if err != nil {
				return
			}
			this_.onMessage(this_, &packet{Index: data[0], Msg: msg})
		}
	})
}
