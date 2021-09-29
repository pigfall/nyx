package server

import(
	"context"
	"encoding/json"
	ws "github.com/gorilla/websocket"
	"github.com/pigfall/tzzGoUtil/net"
	"github.com/pigfall/tzzGoUtil/log"
	"github.com/pigfall/nyx/proto"
)

type connCtrl struct {
	conns map[string]*ws.Conn
	tunIp *net.IpWithMask
	ipPoolIfce ipPoolIfce
	rawLogger log.Logger_Log
}

func newConnCtrl(ipPoolIfce ipPoolIfce,rawLogger log.Logger_Log)*connCtrl {
	return &connCtrl{
		conns :make(map[string]*ws.Conn),
		ipPoolIfce:ipPoolIfce,
		rawLogger :rawLogger,
	}
}


func (this *connCtrl) Serve(
	ctx context.Context,
	conn *ws.Conn,
	tunIfce net.TunIfce,
) (error) {
	logger := log.NewHelper("Serve",this.rawLogger,log.LevelDebug)
	clientVPNIpNet,err := this.ipPoolIfce.GetIpNet()
	if err != nil {
		logger.Error(err)
		return err
	}
	defer this.ipPoolIfce.Release(clientVPNIpNet)
	var clientVPNIpNetStr = clientVPNIpNet.String()
	this.conns[clientVPNIpNetStr ] = conn
	defer delete(this.conns,clientVPNIpNetStr )

	err = connToTunIfce(ctx,this.rawLogger,conn,tunIfce)
	if err != nil {
		logger.Error("connToTunIfce return err %v",err)
		return err
	}
	return nil
}

func connToTunIfce(ctx context.Context,rawLogger log.Logger_Log,conn *ws.Conn,tunIfce net.TunIfce)error{
	logger := log.NewHelper("connToTunIfce",rawLogger,log.LevelDebug)
	for{
		msgType,msgBytes,err := conn.ReadMessage()
		if err != nil{
			logger.Error(err)
			return err
		}
		if msgType == ws.BinaryMessage { // proxy ip packet
			_,err = tunIfce.Write(msgBytes)
			if err != nil{
				logger.Error("write ip packet to tun ifce failed ",err)
				continue
			}
		}else{
			// handle custome proto
			msg := &proto.Msg{}
			err := json.Unmarshal(msgBytes,msg)
			if err != nil {
				logger.Error("parse custome msg failed %w",err)
				continue
			}
			//
		}
	}
}

