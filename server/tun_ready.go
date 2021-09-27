package server

import(
	"context"

	"github.com/pigfall/tzzGoUtil/net"
	"github.com/pigfall/tzzGoUtil/log"
)

// * create tun ifce
// * find the suitable ip , then  assign to the tun ifce
// * enable tun ifce
func tunReady(ctx context.Context,logger log.LoggerLite)error{
	// { collect all net ifce
	allIpV4Addrs,err := net.ListIpV4Addrs()
	if err != nil{
		logger.Error(err)
		return err
	}
	for {

	}
	// }
}
