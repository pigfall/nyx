package client

import(
	"context"
	"fmt"
	"os"
	stdnet "net"
	"github.com/pigfall/tzzGoUtil/net"
	"github.com/pigfall/tzzGoUtil/async"
	"github.com/pigfall/tzzGoUtil/log"
  yy "github.com/pigfall/yingying"
	water_wrap "github.com/pigfall/tzzGoUtil/net/water_tun_wrap"
)


func readyTun(
	ctx context.Context,
	logger log.LoggerLite,
	tunIp *net.IpWithMask,
	asyncCtrl *async.Ctrl,
	tp yy.Transport,
	serverIp stdnet.IP,
)(tunIfce net.TunIfce){
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
	tunIfce = tun
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
	targetA,err := net.FromIpSlashMask("0.0.0.0/1")
	if err != nil{
		panic(err)
	}
	targetB,err := net.FromIpSlashMask("128.0.0.0/1")
	if err != nil{
		panic(err)
	}
	err = net.AddRouteIpNet(targetA,tunIfce.Name(),nil)
	if err != nil{
		logger.Error("Set route table failed %v",err)
		os.Exit(1)
	}
	err = net.AddRouteIpNet(targetB,tunIfce.Name(),nil)
	if err != nil{
		logger.Error("Set route table failed %v",err)
		os.Exit(1)
	}
	logger.Info("Setted route table")
	return tunIfce
}
