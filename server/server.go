package server

import(
	"fmt"
	"context"
	"github.com/pigfall/tzzGoUtil/log"
	"github.com/pigfall/tzzGoUtil/net"
	yy "github.com/pigfall/yingying"
	"github.com/pigfall/tzzGoUtil/async"
)



// serve once
// ready tun interface
// ready connection
func Serve(
	ctx context.Context,
	rawLogger log.Logger_Log,
	cfg *ServeCfg,
	transportServer  yy.TransportServer,
)error{
	asyncCtrl := &async.Ctrl{}
	logger := log.NewHelper("Serve",rawLogger,log.LevelDebug)
	// { ready tun interface
	tunIfce,tunIp, err := tunReady(ctx,logger)
	if err != nil {
		err = fmt.Errorf("准备 tun ifce 环境失败: %w",err)
		logger.Error(err)
		return err
	}
	defer tunIfce.Close()
	asyncCtrl.AppendCancelFuncs(func(){tunIfce.Close()})
	asyncCtrl.OnRoutineQuit(func(){asyncCtrl.Cancel()})
	defer asyncCtrl.Wait()
	defer asyncCtrl.Cancel()
	ipPool,err :=net.NewIpPool(
		tunIp.BaseIpNet(),
		[]*net.IpWithMask{
			tunIp,
		},
	)
	if err != nil{
		panic(err)
	}
	connCtrl := newConnCtrl(ipPool,rawLogger)
	asyncCtrl.AsyncDo(
		ctx,
		func(ctx context.Context){
			var buf = make([]byte,1024*4)
			for {
				n,err := tunIfce.Read(buf)
				if err != nil{
					logger.Error(err)
					return
				}
				for _,conn := range  connCtrl.GetConns(){
					// TODO
					err =conn.WriteIpPacket(buf[:n])
					if err != nil{
						logger.Error(err)
					}
				}
			}
		},
	)
	// }

	// { socket listen
	err = transportServer.Serve(ctx,logger,connCtrl,tunIfce,tunIp)
	defer asyncCtrl.Cancel()
	if err != nil{
		logger.Error(err)
		return err
	}
	// }
	return nil
}


