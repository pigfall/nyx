package client

import (
	"context"
	"fmt"
	stdnet "net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/pigfall/tzzGoUtil/async"
	"github.com/pigfall/tzzGoUtil/log"
	"github.com/pigfall/tzzGoUtil/net"
	yy "github.com/pigfall/yingying"
	// water_wrap "github.com/pigfall/tzzGoUtil/net/water_tun_wrap"
)

func readyTun(
	ctx context.Context,
	logger log.LoggerLite,
	tunIp *net.IpWithMask,
	asyncCtrl *async.Ctrl,
	tp yy.Transport,
	serverIp stdnet.IP,
) (tunIfce net.TunIfce, err error) {
	logger.Info("Creating tun ifce ")
	tun, err := NewTun()
	if err != nil {
		err = fmt.Errorf("Create tun ifce failed %v", err)
		logger.Error(err)
		return nil, err
	}
	logger.Info("Created tun ifce")
	tunName,err := tun.Name()
	if err != nil{
		logger.Error(err)
		return nil,err
	}
	netIfce,err := net.FindIfceByName(tunName)
	if err != nil {
		logger.Error(err)
		return nil,err
	}
	err = tun.SetIp(tunIp.String())
	if err != nil {
		logger.Error("Set tun ifce set failed %w", err)
		return nil, err
	}
	tunIfce = tun
	asyncCtrl.AppendCancelFuncs(func() { tun.Close() })
	asyncCtrl.AppendCancelFuncs(func() {
		net.DelRoute(serverIp)
	})
	asyncCtrl.AsyncDo(
		ctx,
		func(ctx context.Context) {
			<-ctx.Done()
			tun.Close()
		},
	)
	asyncCtrl.AsyncDo(
		ctx,
		func(ctx context.Context) {
			var buf = make([]byte, 1024*4)
			for {
				n, err := tun.Read(buf)
				var data = buf[:n]
				if err != nil {
					logger.Error(err)
					return
				}
				err = tp.WriteIpPacket(data)
				if err != nil {
					logger.Error(err)
				}
			}
		},
	)
	logger.Info("Create tun success")
	logger.Info("Setting route table")
	tunIfceName, err := tun.Name()
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	err = readyTunRoute(logger, serverIp, tunIfceName)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	logger.Info("Setted route table")
	logger.Info("<<<< Setting dns server")
	err = readyTUNDNS(logger,stdnet.ParseIP("8.8.8.8"),netIfce.Index)
	if err != nil{
		logger.Error(">>>> Failed to set dns server")
	}
	logger.Info(">>>> Setted dns server")
	return tunIfce, nil
}

func isDNSReq(ipPacket []byte, logger log.LoggerLite) *dnsToQuery {
	packet := gopacket.NewPacket(ipPacket, layers.LayerTypeIPv4, gopacket.Default)
	if errLayer := packet.ErrorLayer(); errLayer != nil {
		packet = gopacket.NewPacket(ipPacket, layers.LayerTypeIPv6, gopacket.Default)
		errLayer = packet.ErrorLayer()
		if errLayer != nil {
			logger.Error("parse ip failed")
		}

		// logger.Error("parse ipv4 packet failed: ", errLayer)
		return nil
	}
	dnsLayer := packet.Layer(layers.LayerTypeDNS)
	if dnsLayer != nil {
		logger.Debug("is dns req")
		return &dnsToQuery{
			packet:   packet,
			layerDNS: (dnsLayer).(*layers.DNS),
		}
	}

	return nil
}
