package consulclient

import (
	"net"
	"strconv"
	"time"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc/naming"

	"github.com/njmdk/common/logger"
	"github.com/njmdk/common/utils"
)

// NewConsulResolver consul resolver
func NewConsulResolver(address string, service string, log *logger.Logger) naming.Resolver {
	return &consulResolver{
		address: address,
		service: service,
		log:     log,
	}
}

type consulResolver struct {
	address string
	service string
	log     *logger.Logger
}

// Resolve implement
func (r *consulResolver) Resolve(target string) (naming.Watcher, error) {
	config := api.DefaultConfig()
	config.Address = r.address
	config.TokenFile = "./token/token.server"
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	ips, err := GetIPS(client)
	if err != nil {
		return nil, err
	}
	ipS := make([]string, 0, 2)
	ipS = append(ipS, ips.Lan)
	ipS = append(ipS, ips.Wan)
	return &consulWatcher{
		client:  client,
		service: r.service,
		addrS:   map[string]struct{}{},
		log:     r.log,
		SelfIpS: ipS,
	}, nil
}

type consulWatcher struct {
	client    *api.Client
	service   string
	addrS     map[string]struct{}
	lastIndex uint64
	log       *logger.Logger
	SelfIpS   []string
}

func (w *consulWatcher) Next() ([]*naming.Update, error) {
	w.log.Info("start consulWatcher next", zap.String("w.service", w.service))
	for {
		services, metaInfo, err := w.client.Health().Service(w.service, "", true, &api.QueryOptions{
			WaitIndex: w.lastIndex, // 同步点，这个调用将一直阻塞，直到有新的更新
		})
		if err != nil { // 这里不能返回错误，不然会导致consul watcher 停止工作
			w.log.Warn("error retrieving instances from Consul", zap.Error(err))
			time.Sleep(time.Millisecond * 200)
			w.lastIndex = 0
			continue
		}
		w.lastIndex = metaInfo.LastIndex

		addrS := map[string]struct{}{}
		for _, service := range services {
			addr := utils.CheckIsSameIP(net.JoinHostPort(service.Service.Address, strconv.Itoa(service.Service.Port)), w.SelfIpS...)
			addrS[addr] = struct{}{}
		}

		var updates []*naming.Update
		for addr := range w.addrS {
			if _, ok := addrS[addr]; !ok {
				updates = append(updates, &naming.Update{Op: naming.Delete, Addr: addr})
			}
		}

		for addr := range addrS {
			if _, ok := w.addrS[addr]; !ok {
				updates = append(updates, &naming.Update{Op: naming.Add, Addr: addr})
			}
		}

		if len(updates) != 0 {
			w.addrS = addrS
			w.log.Info("获取到grpc地址", zap.Any("addr", addrS))
			return updates, nil
		}
	}
}

func (w *consulWatcher) Close() {
	// nothing to do
}
