package client

import(
	stdnet "net"
	"fmt"
	"github.com/pigfall/tzzGoUtil/net"
	"github.com/pigfall/tzzGoUtil/log"
)


func readyTunRoute(logger log.LoggerLite,serverIp stdnet.IP,tunIfceName string)(error){
	defaultRoute,err := net.GetDefaultRouteRule()
	if err != nil{
		err = fmt.Errorf("Get default route from kernel failed: %w",err)
		logger.Error(err)
		return err
	}
	net.DelRoute(serverIp)
	err = net.AddRoute(serverIp,defaultRoute.DevName,defaultRoute.Via)
	if err != nil{
		logger.Error(err)
		return err
	}
	targetA,err := net.FromIpSlashMask("0.0.0.0/1")
	if err != nil{
		logger.Error(err)
		return err
	}
	targetB,err := net.FromIpSlashMask("128.0.0.0/1")
	if err != nil{
		logger.Error(err)
		return err
	}
	err = net.AddRouteIpNet(targetA,tunIfceName,nil)
	if err != nil{
		logger.Error("Set route table failed %v",err)
		return err
	}
	err = net.AddRouteIpNet(targetB,tunIfceName,nil)
	if err != nil{
		logger.Error("Set route table failed %v",err)
		return err
	}

	return nil
}
