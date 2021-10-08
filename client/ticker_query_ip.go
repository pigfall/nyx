package client

import(
	"fmt"
	ctxLib "github.com/pigfall/tzzGoUtil/ctx"
	"github.com/pigfall/tzzGoUtil/log"
 yy "github.com/pigfall/yingying"
 "github.com/pigfall/yingying/proto"
	"context"
	"time"
)



func tickerQueryIp(ctx context.Context,tp yy.Transport,logger log.LoggerLite){
	ctxLib.TickerDo(ctx,time.Second*2,func()error{
		err := tp.WriteMsg(
			&proto.Msg{
				Id:proto.ID_C2S_QUERY_IP,
			},
			nil,
		)
		if err != nil{
			err = fmt.Errorf("Query ip failed %w",err)
			logger.Error(err)
			return nil // continue query
		}
		logger.Info("Suc send query ip request")
		return nil
	})
	logger.Info("TickerQueryIp quit")
}
