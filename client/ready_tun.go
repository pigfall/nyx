package client

import(
	"context"
	"fmt"
	"os"
	"github.com/pigfall/tzzGoUtil/net"
	"github.com/pigfall/tzzGoUtil/async"
	"github.com/pigfall/tzzGoUtil/log"
  yy "github.com/pigfall/yingying"
	water_wrap "github.com/pigfall/tzzGoUtil/net/water_tun_wrap"
)


func readyTun(ctx context.Context,logger log.LoggerLite,tunIp *net.IpWithMask,tunIfce *net.TunIfce,asyncCtrl *async.Ctrl,tp yy.Transport){
	tun,err := water_wrap.NewTun()
	if err != nil{
		err = fmt.Errorf("Create tun ifce failed %v",err)
		logger.Error(err)
		os.Exit(1)
	}
	err = tun.SetIp(tunIp.String())
	if err != nil{
		logger.Error("Set tun ifce set failed %w",err)
		os.Exit(1)
	}
	*tunIfce = tun
	asyncCtrl.AppendCancelFuncs(func(){tun.Close()})
	asyncCtrl.AsyncDo(
		ctx,
		func(ctx context.Context){
			var buf = make([]byte,1024*4)
			for {
				n,err := tun.Read(buf)
				if err != nil{
					logger.Error(err)
					return
				}
				err = tp.WriteIpPacket(buf[:n])
				if err != nil{
					logger.Error(err)
				}
			}
		},
	)
	logger.Info("Create tun success")
	logger.Info("Setting route table")
}
