package server

import(
	"fmt"
	"context"

	"github.com/pigfall/tzzGoUtil/net"
	water_wrap "github.com/pigfall/tzzGoUtil/net/water_tun_wrap"
	"github.com/pigfall/tzzGoUtil/log"
)

// * create tun ifce
// * find the suitable ip , then  assign to the tun ifce
// * enable tun ifce
func tunReady(ctx context.Context,logger log.LoggerLite)(tunIfce net.TunIfce,tunIp *net.IpWithMask,err error){
	// { collect all net ip, and select the ip which not conflict with been used
	allIpV4Addrs,err := net.ListIpV4Addrs()
	if err != nil{
		logger.Error(err)
		return nil, nil, err
	}
	tunCidr,err := findSuitableIp(allIpV4Addrs)
	if err != nil{
		logger.Error(err)
		return nil, nil, err
	}
	logger.Info("tun cidr ",tunCidr.String())
	// }

	// { create tun ifce 
	tunIfce,err = water_wrap.NewTun()
	if err != nil{
		logger.Error(err)
		return  nil, nil, err
	}
	// }

	// { set ip to tun ifce and enable it
	err = tunIfce.SetIp(tunCidr.String())
	if err != nil{
		err = fmt.Errorf("Set ip to tun interface failed: %w",err)
		logger.Error(err)
		return nil, nil, err
	}
	// }
	tunIp = tunCidr
	return tunIfce, tunIp,nil
}

func findSuitableIp(allIpV4IsUsed []net.IpWithMask)(*net.IpWithMask,error){
	var subnet = 8
	var subnet2 = 0
	encodeIpNet := func(subNet2,subNet int)string{
		return fmt.Sprintf("10.%d.%d.1/16",subNet2,subNet)
	}
	OUT:
	for{
		subnet2++
		if subnet2 >=255{
				return nil,fmt.Errorf("Over then 255 , not found unconflict ip to tun ifce")
		}
		subnet = 1
		for{
			retIp ,err := net.FromIpSlashMask(encodeIpNet(subnet2,subnet))
			if err != nil{
				return nil,err
			}
			if net.IpSubnetCoincideOrCoinCided(retIp,allIpV4IsUsed){
				subnet++
			}else{
				return retIp,nil
			}
			if subnet >= 255{
				continue OUT
			}
		}
	}



}
