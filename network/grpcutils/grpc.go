package grpcutils

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"


	consulclient "github.com/njmdk/common/consul_client"
	"github.com/njmdk/common/logger"
	"github.com/njmdk/common/timer"
)

func NewGRPCClient(addr string, serverName string, log *logger.Logger, isDebug bool) (*grpc.ClientConn, error) {
	unaryInterceptor := func() grpc.UnaryClientInterceptor {
		retry := grpc_retry.UnaryClientInterceptor(grpc_retry.WithBackoff(grpc_retry.BackoffLinear(time.Millisecond*100)),
			grpc_retry.WithMax(2),
			grpc_retry.WithPerRetryTimeout(time.Second*10),
			grpc_retry.WithCodes(4, 8, 10))
		return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			if isDebug {
				reqProto := req.(proto.Message)
				replyProto := reply.(proto.Message)
				start := timer.Now()
				defer log.Info("grpc client request", zap.String("method", method),
					zap.String("req name", proto.MessageName(reqProto)), zap.String("reply name", proto.MessageName(replyProto)),
					zap.Any("req data", reqProto), zap.Any("reply data", replyProto),
					zap.Duration("costs", timer.Now().Sub(start)))

			}
			return retry(ctx, method, req, reply, cc, invoker, opts...)
		}
	}

	cp := keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}
	client, err := grpc.Dial("",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(unaryInterceptor()),
		grpc.WithKeepaliveParams(cp),
		grpc.WithBalancer(grpc.RoundRobin(consulclient.NewConsulResolver(addr, serverName, log))),
	)

	return client, err
}

func NewGRPCServer(consulClient *consulclient.ConsulClient, serverName string, listenPort int, log *logger.Logger, isDebug bool) (*grpc.Server, error) {
	unaryInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := timer.Now()
		defer func() {
			reqProto := req.(proto.Message)
			e := recover()
			if e != nil {
				log.Error("grpc recv response panic", zap.Any("panic info", e), zap.String("method", info.FullMethod),
					zap.String("req name", proto.MessageName(reqProto)),
					zap.Any("req data", reqProto), zap.Duration("costs", timer.Now().Sub(start)))
			} else {
				if err != nil {
					log.Warn("grpc recv response", zap.Error(err), zap.String("method", info.FullMethod),
						zap.String("req name", proto.MessageName(reqProto)),
						zap.Any("req data", reqProto), zap.Duration("costs", timer.Now().Sub(start)))
				} else {
					respProto := resp.(proto.Message)
					if isDebug {
						if info.FullMethod != "/grpc.health.v1.Health/Check" {
							log.Debug("grpc recv response", zap.String("method", info.FullMethod),
								zap.String("req name", proto.MessageName(reqProto)),
								zap.Any("req data", reqProto),
								zap.String("resp name", proto.MessageName(respProto)),
								zap.Any("resp data", respProto), zap.Duration("costs", timer.Now().Sub(start)))
						}
					}
				}
			}
		}()
		resp, err = handler(ctx, req)
		return
	}
	server := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     time.Second * 15,
		MaxConnectionAge:      time.Minute * 30,
		MaxConnectionAgeGrace: time.Second * 5,
		Time:                  time.Second * 10,
		Timeout:               time.Second * 5,
	}), grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		MinTime:             time.Second * 5,
		PermitWithoutStream: true,
	}), grpc.ConnectionTimeout(time.Second*2), grpc.UnaryInterceptor(unaryInterceptor))
	err := consulClient.RegisterGRPCService(serverName, listenPort, server)
	if err != nil {
		return nil, err
	}
	return server, nil
}
