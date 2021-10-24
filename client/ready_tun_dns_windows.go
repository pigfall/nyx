package client

import(
	"golang.org/x/sys/windows"
	"fmt"
		
	gonet "net"
	"golang.zx2c4.com/wireguard/windows/tunnel/winipcfg"
	"github.com/pigfall/tzzGoUtil/log"
	"github.com/pigfall/tzzGoUtil/net"
)

func readyTUNDNS(logger log.LoggerLite,dnsServer gonet.IP,devIndex int) (error){
	luid,err := winipcfg.LUIDFromIndex(uint32(devIndex))
	if err != nil{
		err = fmt.Errorf("LUID from dev index failed: %w",err)
		logger.Error(err)
		return err
	}
	var addressFamily int
	if net.IsIpv4(dnsServer){
		addressFamily = windows.AF_INET
	}else if net.IsIpv6(dnsServer){
		addressFamily = windows.AF_INET6
	}else{
		err = fmt.Errorf("UnKnown ip type %v",dnsServer)
		logger.Error(err)
		panic(err)
	}
	err =luid.SetDNS(winipcfg.AddressFamily(addressFamily),[]gonet.IP{dnsServer},nil)
	if err != nil{
		err = fmt.Errorf("Windows Set DNS %v failed to net interface %v failed: %w",dnsServer,devIndex,err)
		logger.Error(err)
		return err
	}

	return nil
}
