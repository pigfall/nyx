package server

import(
		"net"
		pg_net "github.com/pigfall/tzzGoUtil/net"
)


type ipPoolIfce interface{
	GetIpNet()(*pg_net.IpWithMask,error)
	SetIpCIDR(ip *pg_net.IpWithMask)
	MarkIpAsUsed(ip net.IP)
}

