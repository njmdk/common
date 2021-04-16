package consulclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/njmdk/common/db"
	"github.com/njmdk/common/logger"
	"github.com/njmdk/common/redisclient"
	aliyunsms "github.com/njmdk/common/sms/aliyun_sms"
	"github.com/njmdk/common/utils"
)

type ConsulClient struct {
	client     *consulApi.Client
	serverName string
	log        *logger.Logger
	ip         string
}

func (this_ *ConsulClient) GetClient() *consulApi.Client {
	return this_.client
}

func (this_ *ConsulClient) GetConsulIP() string {
	return this_.ip
}

var MysqlNotFound = errors.New("not found mysql config")

func (this_ *ConsulClient) GetMysql(name string) (*db.Config, error) {
	kv, _, err := this_.client.KV().Get(name, nil)
	if err != nil {
		return nil, err
	}
	if kv == nil {
		return nil, MysqlNotFound
	}
	cfg := &db.Config{}
	err = json.Unmarshal(kv.Value, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, err
}

func (this_ *ConsulClient) GetRedis(name string) (*redisclient.Config, error) {
	kv, _, err := this_.client.KV().Get(name, nil)
	if err != nil {
		return nil, err
	}
	if kv == nil {
		return nil, fmt.Errorf("can not found redis key:%s", name)
	}
	cfg := &redisclient.Config{}
	err = json.Unmarshal(kv.Value, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, err
}

//func (this_ *ConsulClient) GetSSDB(name string) (*redisclient.SSDBConfig, error) {
//	kv, _, err := this_.client.KV().Get(name, nil)
//	if err != nil {
//		return nil, err
//	}
//	if kv == nil {
//		return nil, fmt.Errorf("can not found redis key:%s", name)
//	}
//	cfg := &redisclient.SSDBConfig{}
//	err = json.Unmarshal(kv.Value, cfg)
//	if err != nil {
//		return nil, err
//	}
//	return cfg, err
//}

func (this_ *ConsulClient) GetServices(serverName string) ([]string, error) {
	services, _, err := this_.client.Health().Service(serverName, "", true, nil)
	if err != nil {
		return nil, err
	}
	var addrS []string
	for _, v := range services {
		addrS = append(addrS, net.JoinHostPort(v.Service.Address, strconv.Itoa(v.Service.Port)))
	}
	return addrS, nil
}

func (this_ *ConsulClient) GetSMSInfo(key string) (*aliyunsms.SmsStruct, error) {
	kv, _, err := this_.client.KV().Get(key, nil)
	if err != nil {
		return nil, err
	}
	if kv == nil {
		return nil, fmt.Errorf("can not found sms key:%s", key)
	}
	cfg := &aliyunsms.SmsStruct{}
	err = json.Unmarshal(kv.Value, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, err
}

var ErrorNotFound = errors.New("not found key")

func (this_ *ConsulClient) GetKey(key string) ([]byte, error) {
	kv, _, err := this_.client.KV().Get(key, nil)
	if err != nil {
		return nil, err
	}
	if kv == nil {
		return nil, ErrorNotFound
	}
	return kv.Value, nil
}

func (this_ *ConsulClient) GetServicesIPs(serverName string) ([]*Ips, error) {
	services, _, err := this_.client.Health().Service(serverName, "", true, nil)
	if err != nil {
		return nil, err
	}
	var addrS []*Ips
	for _, v := range services {
		lan := net.JoinHostPort(v.Service.Address, strconv.Itoa(v.Service.Port))
		var wan string
		if v.Service.TaggedAddresses != nil {
			wanServer := v.Service.TaggedAddresses["wan"]
			if wanServer.Address != "" && wanServer.Port > 0 {
				wan = net.JoinHostPort(wanServer.Address, strconv.Itoa(wanServer.Port))
			}
		}
		if lan == "" {
			lan = wan
		}
		if wan == "" {
			wan = lan
		}
		addrS = append(addrS, &Ips{
			Lan: lan,
			Wan: wan,
		})
	}
	return addrS, nil
}

func NewConsulClient(addr string, serverName string, log *logger.Logger) (*ConsulClient, error) {
	cfg := consulApi.DefaultConfig()
	cfg.Address = addr

	cfg.TokenFile = "./token/token.server"
	c, err := consulApi.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	ip, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	cc := &ConsulClient{
		client:     c,
		serverName: serverName,
		log:        log,
		ip:         ip,
	}
	return cc, nil
}

func (this_ *ConsulClient) HealthCheckEcho(e *echo.Echo) string {
	url := "/consul/echo/check"
	e.Any(url, func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "consulCheck")
	})
	return url
}

func (this_ *ConsulClient) HealthCheckDefault(mux http.ServeMux) string {
	url := "/consul/defalut/check"
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("consulCheck"))
	})
	return url
}

type Ips struct {
	Lan string `json:"lan"`
	Wan string `json:"wan"`
}

// Check 实现健康检查接口，这里直接返回健康状态，这里也可以有更复杂的健康检查策略，比如根据服务器负载来返回
func (this_ *ConsulClient) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	if req.Service != this_.serverName {
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_UNKNOWN,
		}, nil
	}

	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}
func (this_ *ConsulClient) Watch(req *grpc_health_v1.HealthCheckRequest, watch grpc_health_v1.Health_WatchServer) error {
	if req.Service != this_.serverName {
		return watch.Send(&grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_UNKNOWN,
		})
	}
	return watch.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
}

func (this_ *ConsulClient) RegisterGRPCService(serverName string, port int, server *grpc.Server) error {
	//this_.RegisterService(serverName, port, &consulApi.AgentServiceCheck{
	//	GRPC:                           fmt.Sprintf("http://%s:%d%s", serverName, port, url),
	//	Timeout:                        "2s",
	//	Interval:                       "2s",
	//	DeregisterCriticalServiceAfter: "2s",
	//})

	grpc_health_v1.RegisterHealthServer(server, this_)
	kv, _, err := this_.client.KV().Get("iprequest", nil)
	if err != nil {
		return err
	}
	ips := &Ips{}
	err = json.Unmarshal(kv.Value, ips)
	if err != nil {
		return err
	}
	if ips.Lan == "" && ips.Wan == "" {
		return errors.New("please add consul kv of iprequest")
	}
	var asc consulApi.AgentServiceChecks
	var lanIP, wanIP string
	serverID := ""
	if ips.Lan != "" {
		lanIP, err = this_.GetSelfIp(ips.Lan)
		if err != nil {
			return err
		}
		serverID = fmt.Sprintf("%s_%s_%d", serverName, lanIP, port)
		lanCheck := &consulApi.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d/%s", lanIP, port, this_.serverName),
			Timeout:                        "1s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "1m",
		}
		asc = append(asc, lanCheck)
	}
	if ips.Wan != "" {
		wanIP, err = this_.GetSelfIp(ips.Wan)
		if err != nil {
			return err
		}
		if serverID == "" {
			serverID = fmt.Sprintf("%s_%s_%d", serverName, wanIP, port)
		}
		wanCheck := &consulApi.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d/%s", wanIP, port, this_.serverName),
			Timeout:                        "2s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "1m",
		}
		asc = append(asc, wanCheck)
	}

	register := new(consulApi.AgentServiceRegistration)
	register.ID = serverID
	register.Name = serverName

	if lanIP == "" {
		register.Address = wanIP
	} else {
		register.Address = lanIP
		if wanIP != "" {
			if register.TaggedAddresses == nil {
				register.TaggedAddresses = map[string]consulApi.ServiceAddress{}
			}
			register.TaggedAddresses["wan"] = consulApi.ServiceAddress{
				Address: wanIP,
				Port:    port,
			}
		}
	}

	register.Port = port
	register.Checks = asc
	this_.timer(time.Second*2, func() {
		err = this_.client.Agent().ServiceRegister(register)
		if err != nil {
			this_.log.Error("server register error", zap.Error(err), zap.Any("data", register))
		}
	})
	return nil
}

func (this_ *ConsulClient) RegisterHTTPService(serverName string, port int, url string) error {
	kv, _, err := this_.client.KV().Get("iprequest", nil)
	if err != nil {
		return err
	}
	if kv == nil {
		return errors.New("can not found key:iprequest")
	}
	ips := &Ips{}
	err = json.Unmarshal(kv.Value, ips)
	if err != nil {
		return err
	}
	if ips.Lan == "" && ips.Wan == "" {
		return errors.New("please add consul kv of iprequest")
	}
	var asc consulApi.AgentServiceChecks
	var lanIP, wanIP string
	serverID := ""
	if ips.Lan != "" {
		lanIP, err = this_.GetSelfIp(ips.Lan)
		if err != nil {
			return err
		}

		serverID = fmt.Sprintf("%s_%s_%d", serverName, lanIP, port)
		lanCheck := &consulApi.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d%s", lanIP, port, url),
			Timeout:                        "2s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "1m",
		}
		asc = append(asc, lanCheck)
	}
	if ips.Wan != "" {
		wanIP, err = this_.GetSelfIp(ips.Wan)
		if err != nil {
			return err
		}
		if serverID == "" {
			serverID = fmt.Sprintf("%s_%s_%d", serverName, wanIP, port)
		}
		wanCheck := &consulApi.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d%s", wanIP, port, url),
			Timeout:                        "2s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "1m",
		}
		asc = append(asc, wanCheck)
	}

	register := new(consulApi.AgentServiceRegistration)
	register.ID = serverID
	register.Name = serverName

	if lanIP == "" {
		register.Address = wanIP
	} else {
		register.Address = lanIP
		if wanIP != "" {
			if register.TaggedAddresses == nil {
				register.TaggedAddresses = map[string]consulApi.ServiceAddress{}
			}
			register.TaggedAddresses["wan"] = consulApi.ServiceAddress{
				Address: wanIP,
				Port:    port,
			}
		}
	}

	register.Port = port
	register.Checks = asc
	this_.timer(time.Second*2, func() {
		err = this_.client.Agent().ServiceRegister(register)
		if err != nil {
			this_.log.Error("server register error", zap.Error(err), zap.Any("data", register))
		}
	})
	return nil
}

func (this_ *ConsulClient) timer(duration time.Duration, f func()) {
	panicFunc := func(i interface{}) {
		this_.log.Error("consul register safego panic", zap.Any("panic info", i))
	}
	utils.SafeGO(panicFunc, func() {
		timer := time.NewTimer(duration)
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				func() {
					defer utils.Recover(panicFunc)
					f()
				}()
				timer.Reset(duration)
			}
		}
	})
}

func (this_ *ConsulClient) RegisterTCPService(serverName string, port int) error {
	kv, _, err := this_.client.KV().Get("iprequest", nil)
	if err != nil {
		return err
	}
	if kv == nil {
		return errors.New("can not found key:iprequest")
	}
	ips := &Ips{}
	err = json.Unmarshal(kv.Value, ips)
	if err != nil {
		return err
	}
	if ips.Lan == "" && ips.Wan == "" {
		return errors.New("please add consul kv of iprequest")
	}
	var asc consulApi.AgentServiceChecks
	var lanIP, wanIP string
	serverID := ""
	if ips.Lan != "" {
		lanIP, err = this_.GetSelfIp(ips.Lan)
		if err != nil {
			return err
		}

		serverID = fmt.Sprintf("%s_%s_%d", serverName, lanIP, port)
		lanCheck := &consulApi.AgentServiceCheck{
			TCP:                            fmt.Sprintf("%s:%d", lanIP, port),
			Timeout:                        "2s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "1m",
		}
		asc = append(asc, lanCheck)
	}
	if ips.Wan != "" {
		wanIP, err = this_.GetSelfIp(ips.Wan)
		if err != nil {
			return err
		}
		if serverID == "" {
			serverID = fmt.Sprintf("%s_%s_%d", serverName, wanIP, port)
		}
		wanCheck := &consulApi.AgentServiceCheck{
			TCP:                            fmt.Sprintf("%s:%d", wanIP, port),
			Timeout:                        "2s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "1m",
		}
		asc = append(asc, wanCheck)
	}

	register := new(consulApi.AgentServiceRegistration)
	register.ID = serverID
	register.Name = serverName

	if lanIP == "" {
		register.Address = wanIP
	} else {
		register.Address = lanIP
		if wanIP != "" {
			if register.TaggedAddresses == nil {
				register.TaggedAddresses = map[string]consulApi.ServiceAddress{}
			}
			register.TaggedAddresses["wan"] = consulApi.ServiceAddress{
				Address: wanIP,
				Port:    port,
			}
		}
	}

	register.Port = port
	register.Checks = asc
	this_.timer(time.Second*2, func() {
		err = this_.client.Agent().ServiceRegister(register)
		if err != nil {
			this_.log.Error("server register error", zap.Error(err), zap.Any("data", register))
		}
	})

	return nil
}

func (this_ *ConsulClient) GetSelfIp(addr string) (string, error) {
	resp, err := http.Get(addr)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	m := map[string]string{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return "", err
	}

	if v, ok := m["ip"]; ok && v != "" {
		return v, nil
	}
	return "", errors.New("not found self ip")
}

func GetSelfIp(addr string) (string, error) {
	resp, err := http.Get(addr)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	m := map[string]string{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return "", err
	}

	if v, ok := m["ip"]; ok && v != "" {
		return v, nil
	}
	return "", errors.New("not found self ip")
}
func (this_ *ConsulClient) GetIPS() (*Ips, error) {
	return GetIPS(this_.client)
}
func GetIPS(client *consulApi.Client) (*Ips, error) {
	kv, _, err := client.KV().Get("iprequest", nil)
	if err != nil {
		return nil, err
	}
	if kv == nil {
		return nil, errors.New("can not found key:iprequest")
	}
	ips := &Ips{}
	err = json.Unmarshal(kv.Value, ips)
	if err != nil {
		return nil, err
	}
	if ips.Lan == "" && ips.Wan == "" {
		return nil, errors.New("please add consul kv of iprequest")
	}
	var lanIP, wanIp string
	if ips.Lan != "" {
		lanIP, err = GetSelfIp(ips.Lan)
		if err != nil {
			return nil, err
		}
	}
	if ips.Wan != "" {
		wanIp, err = GetSelfIp(ips.Wan)
		if err != nil {
			return nil, err
		}
	}
	if wanIp == "" {
		wanIp = lanIP
	}
	if lanIP == "" {
		lanIP = wanIp
	}
	return &Ips{
		Lan: lanIP,
		Wan: wanIp,
	}, nil
}
