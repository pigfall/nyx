package client

import(
	"time"
	"sync"
	"fmt"
		"context"
		"github.com/pigfall/tzzGoUtil/log"
		"github.com/pigfall/tzzGoUtil/async"
		"github.com/pigfall/tzzGoUtil/net"
    ws "github.com/gorilla/websocket"
		"encoding/json"
		yy "github.com/pigfall/yingying"
		"github.com/pigfall/yingying/transport"
		"github.com/pigfall/yingying/proto"
)

func Run(
	ctx context.Context,
	rawLogger log.Logger_Log,
	cfg *RunCfg,
)error{
	wg := sync.WaitGroup{}
	defer wg.Wait()
	ctx,cancel := context.WithCancel(ctx)
	defer cancel()
	// 
	clientState := &clientState{
		cancel:cancel,
	}
	logger := log.NewHelper("ClientRun",rawLogger,log.LevelDebug)
	logger.Info("Connecting Server")
	conn,_,err := ws.DefaultDialer.Dial(cfg.ServerAddr,nil)
	if err != nil{
		logger.Errorf("Failed to connect server %s",cfg.ServerAddr)
		return err
	}
	tp := transport.NewTPWebSocket(conn)
	logger.Info("Connected server")
	async.AsyncDo(
			ctx,
			&wg,
			func (ctx context.Context) {
				logger := log.NewHelper("checkClientIp",rawLogger,log.LevelDebug)
				ticker := time.NewTicker(time.Second*2)
				for{
					select{
					case <-ctx.Done():
						return
					case <-ticker.C:
						logger.Info("Request client ip")
						if !clientState.HasIp(){
							requestClientIp(tp)
						}
					}
				}
			},
	)
	//

	async.AsyncDo(
		ctx,
		&wg,
		func(ctx context.Context){
			clientConnReadHandler(
				ctx,
				tp,
				rawLogger,
			)
		},
	)
	defer conn.Close()
	err = requestClientIp(tp)
	if err != nil {
		err = fmt.Errorf("Failed to request vpn client ip from server %w ",err)
		logger.Error(err)
		return err
	}
	return nil
}

type clientState struct{
	ip  *net.IpWithMask
	l sync.Mutex
	cancel func()
}

func (this *clientState) SetIp(ip *net.IpWithMask) {
	this.l.Lock()
	defer this.l.Unlock()
	if this.ip != nil {
		if this.ip.String() != ip.String(){
			panic(fmt.Errorf("Has been assined a ip %s, but server give a new ip %s",this.ip.String(),ip.String()))
		}
	}else{
		this.ip = ip
		// start create tun ifce
		//
	}

}

func (this *clientState) HasIp() bool{
	this.l.Lock()
	defer this.l.Unlock()
	return this.ip != nil
}



func clientConnReadHandler(ctx context.Context,tp yy.Transport,rawLogger log.Logger_Log)error{
	logger := log.NewHelper("readHandler",rawLogger,log.LevelDebug)
	for {
		msgType,data,err := tp.Read()
		if err != nil{
			logger.Error(err)
			return err
		}
		switch msgType{
		case yy.IpPacket:
			panic("TODO")
		case yy.Proto:
			panic("TODO")
		default:
			panic(fmt.Errorf("Undefined msgType %s",msgType))

		}
	}
}


func handleProtoMsg(tp yy.Transport,logger log.LoggerLite,data []byte,clientState *clientState){
	var msg proto.Msg
	err := json.Unmarshal(data,&msg)
	if err != nil{
		logger.Error(err)
		return
	}
	switch msg.Id{
	case proto.ID_S2C_QUERY_IP:
		var body proto.S2C_ClientVPNIpNet
		err = json.Unmarshal(msg.Body,&body)
		if err != nil{
			logger.Error(err)
			return
		}
		ipNet,err := net.FromIpSlashMask(body.IpNet)
		if err != nil{
			logger.Error(err)
			return
		}
		clientState.SetIp(ipNet)
	}
}

func checkVPNClientIp(){

}

func requestClientIp(tp yy.Transport)error{
	return  tp.WriteMsg(
			&proto.Msg{
				Id:proto.ID_C2S_QUERY_IP,
			},
			nil,
	)
}
