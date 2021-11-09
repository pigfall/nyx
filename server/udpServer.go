package server

import(
	"sync"
	"encoding/json"
	"github.com/pigfall/yingying/proto"
	"net"
	"github.com/pigfall/tzzGoUtil/log"
	tzNet "github.com/pigfall/tzzGoUtil/net"
	"context"
	yy "github.com/pigfall/yingying"
	tp "github.com/pigfall/yingying/transport"
)

type tpServerUDP struct{
	ipToListen net.IP
	port int
	connMgr *connMgr
}

func NewTPServerUDP(ipToListen net.IP,port int,connMgr *connMgr)yy.TransportServer {
	return &tpServerUDP{
		ipToListen:ipToListen,
		port:port,
		connMgr:connMgr,
	}
}

type udpServerTunnelIpMapper struct{
	remoteIpPortToTunnelIp map[tzNet.IpPortFormat]*tzNet.IpWithMask
	tunnelIpToRemoteIpPort map[tzNet.IpFormat]tzNet.IpPortFormat
	ipPool ipPoolIfce
	l sync.Mutex
}

func newUDPServerTunnelIpMapper(ipPool ipPoolIfce)*udpServerTunnelIpMapper{
	return &udpServerTunnelIpMapper{
		ipPool:ipPool,
		remoteIpPortToTunnelIp :make(map[tzNet.IpPortFormat]*tzNet.IpWithMask),
		tunnelIpToRemoteIpPort :make(map[tzNet.IpFormat]tzNet.IpPortFormat),
	}
}

func (this *udpServerTunnelIpMapper)TakeIp(remoteIpPort tzNet.IpPortFormat)(*tzNet.IpWithMask,error){
	this.l.Lock()
	defer this.l.Unlock()
	tunnelIpNet,ok := this.remoteIpPortToTunnelIp[remoteIpPort]
	if ok {
		return tunnelIpNet,nil
	}
	ipNet,err := this.ipPool.Take()
	if err != nil{
		return nil,err
	}
	this.remoteIpPortToTunnelIp[remoteIpPort] = ipNet
	this.tunnelIpToRemoteIpPort[ipNet.IpFormat()] = remoteIpPort

	return ipNet,nil

}

func buildUDPServerIpGetter(remoteIpPort tzNet.IpPortFormat,ipMapper *udpServerTunnelIpMapper)clientTunnelIpGetter{
	return clientTunnelIpGetterFunc(
		func()(tzNet.IpNetFormat,error){
			ipNet,err := ipMapper.TakeIp(remoteIpPort)
			if err != nil{
				return "",err
			}
			return ipNet.ToIpNetFormat(),nil
		},
	)
}


func (this *tpServerUDP)Serve(ctx context.Context,logger log.LoggerLite,connCtrl yy.ConnCtrl,tunIfce tzNet.TunIfce,tunIp *tzNet.IpWithMask)(error){
	udpSock,err := tzNet.UDPListen(this.ipToListen,this.port)
	if err != nil{
		logger.Error(err)
		return err
	}
	ipPool,err :=tzNet.NewIpPool(
		tunIp.BaseIpNet(),
		[]*tzNet.IpWithMask{
			tunIp,
		},
	)
	if err != nil{
		logger.Error(err)
		return err
	}
	tpUDP := tp.NewTransportUDP(udpSock)
	udpServerTunnelIpMapper :=newUDPServerTunnelIpMapper(ipPool)
	var buf = make([]byte,1024*6)
	for {
		n,remote,err := udpSock.ReadFromUDP(buf)
		if err != nil{
			logger.Error(err)
			continue
		}
		// < parse msg
		firstByte := buf[0]
		if firstByte == 0 { // ip packet
			_,err := tunIfce.Write(buf[1:n])
			if err != nil{
				logger.Error(err)
			}
		}else{
			appMsgBytes := buf[1:n]
			// < handle app msg
			// << parse msg
			msg := &proto.Msg{}
			err = json.Unmarshal(appMsgBytes,msg)
			if err != nil{
				logger.Error(err)
				continue
			}

			resMsg := handleAppMsg(
				ctx,
				msg,buildUDPServerIpGetter(
					udpAddrToIpPortFormat(remote), udpServerTunnelIpMapper,
				),
			)
			err = tpUDP.WriteJSON(resMsg)
			if err != nil{
				logger.Error(err)
			}
			// >>
			// >
			// >
		}
	}
}

func udpAddrToIpPortFormat(addr *net.UDPAddr)tzNet.IpPortFormat{
	return tzNet.IpPortFormatFromIpPort(addr.IP,addr.Port)
}
