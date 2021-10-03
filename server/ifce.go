package server

import(
		pg_net "github.com/pigfall/tzzGoUtil/net"
)


type ipPoolIfce interface{
	Take()(*pg_net.IpWithMask,error)
	Release(*pg_net.IpWithMask)
}

