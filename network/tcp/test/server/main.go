package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"

	"github.com/njmdk/common/eventqueue"
	"github.com/njmdk/common/logger"
	"github.com/njmdk/common/network/basepb"
	"github.com/njmdk/common/network/tcp"
	workpool "github.com/njmdk/common/work_pool"
)

/*
type AcceptorInterface interface {
	OnSessionConnected(*Session)
	OnSessionDisConnected(*Session, error)
	OnRPCRequest(*Session, proto.Message) proto.Message
	OnNormalMsg(*Session, proto.Message)
}
*/
type Server struct {
	*tcp.Acceptor
	log *logger.Logger
	*tcp.DefaultDoConcurrent
}

func (this_ *Server) RegisterConcurrentMsg() {
	// this_.DefaultDoConcurrent.RegisterConcurrentMsg(&pbmsg.Account_LoginRequest{}) // 如果收到*pbmsg.Account_LoginRequest 消息，那么将会在协程池里面调用
}

func (this_ *Server) OnSessionConnected(s *tcp.Session) {
	//this_.TestRequest(s)
	this_.Acceptor.OnSessionConnected(s)
}

func (this_ *Server) OnSessionDisConnected(s *tcp.Session, err error) {
	this_.Acceptor.OnSessionDisConnected(s, err)
}

func (this_ *Server) OnRPCRequest(s *tcp.Session, m proto.Message) proto.Message {
	_ = this_.Acceptor.OnRPCRequest(s, m)
	//switch msg := m.(type) {
	//case *pbmsg.Account_LoginRequest: // 如果是收到这个消息，会在协程池里面运行
	//	_ = msg
	//case *pbmsg.Center_GetMysqlRequest: // 因为没注册并发处理，所有收到这个消息是在主逻辑线程里面处理
	//	_ = msg
	//}

	return &basepb.Base_Success{}
}

func (this_ *Server) OnNormalMsg(s *tcp.Session, msg proto.Message) {
	this_.Acceptor.OnNormalMsg(s, msg)
}

func (this_ *Server) TestRequest(s *tcp.Session) {
	s.RequestNoError(&basepb.Base{}, func(m proto.Message) {
		msg, ok := m.(*basepb.Base_Success)
		if ok {
			this_.log.Info("recv response", zap.String("msgID", proto.MessageName(m)), zap.String("msg", msg.String()))
			this_.TestRequest(s)
		}
	})
}

func main() {
	log, err := logger.New("tcp_testserver", ".", zap.DebugLevel, true)
	if err != nil {
		panic(err)
	}

	queue := eventqueue.NewEventQueue(10000, log)
	acceptor, err := tcp.NewAcceptor(":9999", queue, &tcp.ProtoCodeC{}, log, true)
	if err != nil {
		panic(err)
	}
	acceptor.StartAccept()

	s := &Server{
		Acceptor:            acceptor,
		log:                 log,
		DefaultDoConcurrent: &tcp.DefaultDoConcurrent{},
	}
	s.RegisterConcurrentMsg()
	acceptor.SetCallback(s)

	queue.Run(nil, tcp.DispatchMsg(func(event interface{}) {
		switch e := event.(type) {
		case *eventqueue.EventStopped:
			log.Warn("event queue stopped")
		default:
			log.Warn("unknown event", zap.Any("event", e))
		}
	}, workpool.NewWorkPool(500, log)))

	sigShutdown := make(chan os.Signal, 5)
	signal.Notify(sigShutdown, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//signal.Notify(sigShutdown)
	sigCall := <-sigShutdown
	log.Info("recv signal", zap.String("signal", sigCall.String()))
	queue.Stop()
}
