package client

import(
	"encoding/json"
	"fmt"
		"context"
ws "github.com/gorilla/websocket"
	"github.com/pigfall/tzzGoUtil/log"
	"github.com/pigfall/tzzGoUtil/async"
	"github.com/pigfall/tzzGoUtil/net"
 yy "github.com/pigfall/yingying"
  "github.com/pigfall/yingying/transport"
  "github.com/pigfall/yingying/proto"
)


func handleConnData(
	ctx context.Context,
	rawLogger log.Logger_Log,
	conn *ws.Conn,
	asyncCtrl *async.Ctrl,
){
	tp := transport.NewTPWebSocket(conn)
	logger := log.NewHelper("handleConnData",rawLogger,log.LevelDebug)
	logger.Info("Reading msg from server")
	ctxQueryIp,cancelQueryIp:=context.WithCancel(ctx)
	defer cancelQueryIp()
	go func(){
		tickerQueryIp(ctxQueryIp,tp,logger)
	}()
	var clientIp *net.IpWithMask
	var tunIfce  net.TunIfce
	for {
		msgType,data,err := tp.Read()
		if err != nil{
			logger.Error(err)
			return 
		}
		// < request client  ip
		// >
		switch msgType {
		case yy.IpPacket:
			err = handleIpPacket(data,tunIfce)
			if err != nil{
				logger.Error(err)
			}
		case yy.Proto:
			tunIfceTmp := handleConnProto(ctx,data,logger,&clientIp,cancelQueryIp,asyncCtrl,tp)
			if tunIfceTmp != nil{
				tunIfce = tunIfceTmp
			}
		default:
			panic(fmt.Errorf("Undefined msgType %v",msgType))
		}
	}
}

func handleIpPacket(data []byte,tunIfce net.TunIfce)error{
	_,err := tunIfce.Write(data)
	return err
}

func handleConnProto(ctx context.Context,data []byte,logger log.LoggerLite,clientIp **net.IpWithMask,cancelQueryIp func(),asyncCtrl *async.Ctrl,tp yy.Transport)(tunIfce net.TunIfce){
	var msg proto.Msg
	err := json.Unmarshal(data,&msg)
	if err != nil{
		logger.Error(err)
		return 
	}
	var unmarshalBody = json.Unmarshal
	switch msg.Id{
	case proto.ID_S2C_QUERY_IP:
		var body proto.S2C_ClientVPNIpNet
		err := unmarshalBody(msg.Body,&body)
		if err != nil{
			logger.Error(err)
			return
		}
		ipNet,err := net.FromIpSlashMask(body.IpNet)
		if err != nil{
			logger.Error(err)
			return
		}
		if *clientIp == nil{
			logger.Info("Get Client ip ",ipNet.String())
			*clientIp = ipNet
			cancelQueryIp()
			return readyTun(ctx,logger,*clientIp,asyncCtrl,tp)
		}else{
			if (*clientIp).String()!=ipNet.String(){
				panic("Client ip not match")
			}
		}
	default:
		panic(fmt.Errorf("Undefined msg id",msg.Id))
	}
	return nil
}
