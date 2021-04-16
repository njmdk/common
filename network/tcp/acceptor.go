package tcp

import (
	"errors"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"

	"github.com/njmdk/common/eventqueue"
	"github.com/njmdk/common/logger"
	"github.com/njmdk/common/network/basepb"
	"github.com/njmdk/common/network/reuseport"
	"github.com/njmdk/common/utils"
)

type AcceptSession struct {
	*Session
}

func (this_ *AcceptSession) GetSession() *Session {
	return this_.Session
}

type SessionClosed struct {
	Err  error
	Sess *Session
}

func (this_ *SessionClosed) GetSession() *Session {
	return this_.Sess
}

type ErrServerClosed string

func (e *ErrServerClosed) Error() string {
	return "acceptor closed"
}

type Acceptor struct {
	ai         DispatchInterface
	dc         DoConcurrentInterface
	addr       string
	IP         string
	Port       string
	queue      *eventqueue.EventQueue
	close      chan struct{}
	listener   net.Listener
	codeC      CodeC
	logger     *logger.Logger
	dialConn   func(conn net.Conn)
	isDebugLog bool
	sync.Once
}

func (this_ *Acceptor) GetListener() net.Listener {
	return this_.listener
}

func NewAcceptor(addr string, queue *eventqueue.EventQueue, codeC CodeC, logger *logger.Logger, isDebugLog bool) (*Acceptor, error) {
	l, err := reuseport.Listen("tcp4", addr)
	if err != nil {
		return nil, err
	}
	tcpAddr := l.Addr().(*net.TCPAddr)
	port := strconv.Itoa(tcpAddr.Port)
	ip := tcpAddr.IP.String()
	acceptor := &Acceptor{
		addr:       addr,
		queue:      queue,
		codeC:      codeC,
		logger:     logger,
		listener:   l,
		Port:       port,
		IP:         ip,
		close:      make(chan struct{}),
		isDebugLog: isDebugLog,
	}
	acceptor.dialConn = acceptor.startConn
	acceptor.ai = acceptor
	acceptor.dc = acceptor
	return acceptor, nil
}

func (this_ *Acceptor) SetCallback(ai DispatchInterface) {
	this_.ai = ai
}

func (this_ *Acceptor) SetDoConcurrent(dc DoConcurrentInterface) {
	this_.dc = dc
}

func (this_ *Acceptor) DoConcurrent(msgID string) bool {
	return false
}

func (this_ *Acceptor) OnSessionConnected(session *Session) {
	this_.logger.Debug("session connected", zap.String("local", session.LocalAddr()), zap.String("remote", session.remoteAddr))
}

func (this_ *Acceptor) OnSessionDisConnected(session *Session, err error) {
	if this_.isDebugLog {
		this_.logger.Warn("session disconnected", zap.Error(err), zap.String("local", session.LocalAddr()), zap.String("remote", session.remoteAddr))
	}
}

func (this_ *Acceptor) OnRPCRequest(session *Session, msg proto.Message) proto.Message {
	if this_.isDebugLog {
		this_.logger.Debug("session recv rpc request", zap.String("msgID", proto.MessageName(msg)), zap.Any("request", msg), zap.String("local", session.LocalAddr()), zap.String("remote", session.remoteAddr))
	}
	return &basepb.Base_Success{}
}

func (this_ *Acceptor) OnNormalMsg(session *Session, msg proto.Message) {
	if this_.isDebugLog {
		this_.logger.Debug("session normal message", zap.String("msgID", proto.MessageName(msg)), zap.Any("msg", msg), zap.String("local", session.LocalAddr()), zap.String("remote", session.remoteAddr))
	}
}

func (this_ *Acceptor) Close() {
	this_.Do(func() {
		close(this_.close)
		_ = this_.listener.Close()
	})
}

func (this_ *Acceptor) startConn(conn net.Conn) {
	sess := NewSession(conn, this_.queue, true, this_.codeC, this_.logger, time.Second*2, this_.isDebugLog)
	sess.onClose = func(e error) {
		this_.queue.Post(&SessionClosed{
			Err:  e,
			Sess: sess,
		})
	}

	this_.queue.Post(&AcceptSession{
		sess,
	})

	this_.queue.Tick(time.Second*30, func(i time.Time) bool {
		if i.Unix()-sess.lastPongTime > 30 {
			sess.Close(errors.New("ping pong timeout"))
			return false
		}
		return true
	})
	sess.Dispatch = this_.ai
	if this_.dc != nil {
		sess.DoConcurrent = this_.dc
	}
	sess.Start()
}

func (this_ *Acceptor) StartAccept() {

	utils.SafeGO(func(e interface{}) {
		this_.StartAccept()
	}, func() {
		var tempDelay time.Duration // 监听失败时暂停多久重新开始接收
		for {
			conn, e := this_.listener.Accept()
			if e != nil {
				select {
				case <-this_.close:
					return
				default:
				}
				if ne, ok := e.(net.Error); ok && ne.Temporary() {
					if tempDelay == 0 {
						tempDelay = 5 * time.Millisecond
					} else {
						tempDelay *= 2
					}
					if max := 1 * time.Second; tempDelay > max {
						tempDelay = max
					}
					this_.logger.WarnFormat("tcp: Accept error: %v; retrying in %v", e, tempDelay)
					time.Sleep(tempDelay)
					continue
				}
				return
			}
			tempDelay = 0
			this_.dialConn(conn)
		}
	})

}
