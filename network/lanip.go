package network

import "net"

func GetLanIP() (string, error) {
	addrS, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	var lanIP, loopbackIP string
	for _, v := range addrS {
		if ipNet, ok := v.(*net.IPNet); ok {
			if !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					lanIP = ipNet.IP.String()
				}
			} else {
				loopbackIP = ipNet.IP.String()
			}
		}
	}
	if lanIP == "" {
		lanIP = loopbackIP
	}
	return lanIP, nil
}

func SplitHostPort(hostPort string) (host, port string, err error) {
	host, port, err = net.SplitHostPort(host)
	if err != nil {
		return
	}
	if host == "" {
		host = "0.0.0.0"
	}
	return
}

func SplitHostPortToAddr(hostPort string) (string, error) {
	host, port, err := net.SplitHostPort(hostPort)
	if err != nil {
		return "", err
	}
	if host == "" {
		host = "0.0.0.0"
	}
	return net.JoinHostPort(host, port), nil
}
