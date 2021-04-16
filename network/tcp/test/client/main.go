package main

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"

	"github.com/njmdk/common/eventqueue"
	"github.com/njmdk/common/logger"
	"github.com/njmdk/common/network/basepb"
	"github.com/njmdk/common/network/tcp"
)

//type ConnectorInterface interface {
//	OnSessionConnected(*tcp.Connector)
//	OnSessionDisConnected(*tcp.Connector, error)
//	OnRPCRequest(*tcp.Connector, proto.Message) proto.Message
//	OnNormalMsg(*tcp.Connector, proto.Message)
//}

type Client struct {
	*tcp.Connector
	log *logger.Logger
}

func (this_ *Client) OnSessionConnected(s *tcp.Session) {
	this_.Connector.OnSessionConnected(s)
	this_.TestRequest()
}

func (this_ *Client) TestRequest() {
	this_.Connector.RequestNoError(&basepb.Base{}, func(m proto.Message) {
		msg, ok := m.(*basepb.Base_Success)
		if ok {
			this_.log.Info("recv response", zap.String("msgID", proto.MessageName(m)), zap.String("msg", msg.String()))
			this_.TestRequest()
		}
	})
}

func (this_ *Client) OnSessionDisConnected(s *tcp.Session, e error) {
	this_.Connector.OnSessionDisConnected(s, e)
}

func (this_ *Client) OnRPCRequest(s *tcp.Session, msg proto.Message) proto.Message {
	_ = this_.Connector.OnRPCRequest(s, msg)
	return &basepb.Base_Success{}
}

func (this_ *Client) OnNormalMsg(s *tcp.Session, msg proto.Message) {
	this_.Connector.OnRPCRequest(s, msg)
}

func main() {
	log, err := logger.New("tcp_testclient", ".", zap.DebugLevel, true)
	if err != nil {
		panic(err)
	}

	queue := eventqueue.NewEventQueue(10000, log)
	connector := tcp.NewConnector("127.0.0.1:9999", queue, log, "test_client", tcp.WithOpenDebugLog(), tcp.WithReconnect())
	c := &Client{
		Connector: connector,
		log:       log,
	}
	connector.SetCallback(c)
	connector.Connect()
	queue.Run(nil, tcp.DispatchMsg(func(event interface{}) {
		switch e := event.(type) {
		case *eventqueue.EventStopped:
			log.Warn("event queue stopped")
		default:
			log.Warn("unknown event", zap.Any("event", e))
		}
	}))

	sigShutdown := make(chan os.Signal, 5)
	signal.Notify(sigShutdown, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//signal.Notify(sigShutdown)
	sigCall := <-sigShutdown
	log.Info("recv signal", zap.String("signal", sigCall.String()))
	c.Close(errors.New("shutdown"))
	queue.Stop()
}
