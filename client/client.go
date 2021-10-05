package client

import(
	"sync"
	"fmt"
		"context"
		"github.com/pigfall/tzzGoUtil/log"
		"github.com/pigfall/tzzGoUtil/async"
    ws "github.com/gorilla/websocket"
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
		func(ctx context.Context){
			clientConnReadHandler(
				ctx,
				tp,
				rawLogger,
			)
		},
	)
	defer wg.Wait()
	defer conn.Close()
	err = tp.WriteMsg(
			&proto.Msg{
				Id:proto.ID_C2S_QUERY_IP,
			},
			nil,
	)
	if err != nil {
		err = fmt.Errorf("Failed to request vpn client ip from server %w ",err)
		logger.Error(err)
		return err
	}
	return nil
}

func clientConnReadHandler(ctx context.Context,tp yy.Transport,rawLogger log.Logger_Log)error{
	panic("TODO")
}
