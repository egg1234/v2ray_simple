package netLayer

import (
	"errors"
	"net"
	"syscall"

	"golang.org/x/net/route"
)

/*
我们的auto route使用纯命令行方式。

sing-box 使用了另一种系统级别的方式。使用了
golang.org/x/net/route

下面给出一些参考

https://github.com/libp2p/go-netroute

https://github.com/jackpal/gateway/issues/27

https://github.com/GameXG/gonet/blob/master/route/route_windows.go

除了 GetGateway之外，还可以使用更多其他代码
*/
func GetGateway() (ip net.IP, index int, err error) {
	var rib []byte
	rib, err = route.FetchRIB(syscall.AF_INET, syscall.NET_RT_DUMP, 0)
	if err != nil {
		return
	}
	var msgs []route.Message
	msgs, err = route.ParseRIB(syscall.NET_RT_DUMP, rib)
	if err != nil {
		return
	}

	for _, m := range msgs {
		switch m := m.(type) {
		case *route.RouteMessage:
			switch sa := m.Addrs[syscall.RTAX_GATEWAY].(type) {
			case *route.Inet4Addr:
				ip = net.IPv4(sa.IP[0], sa.IP[1], sa.IP[2], sa.IP[3])
			case *route.Inet6Addr:
				ip = make(net.IP, net.IPv6len)
				copy(ip, sa.IP[:])
			}
			index = m.Index

			return

		}
	}
	err = errors.New("no gateway")
	return
}
