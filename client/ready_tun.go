package client

import (
	"context"
	"fmt"
	stdnet "net"

	"github.com/pigfall/tzzGoUtil/async"
	"github.com/pigfall/tzzGoUtil/log"
	"github.com/pigfall/tzzGoUtil/net"
	yy "github.com/pigfall/yingying"
	"github.com/google/gopacket/layers"
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
	var dnsToQueryChan :=make(chan interface{},10)
	asyncCtrl.AsyncDo(
			ctx,
			func(ctx context.Context){
				var toQueryMap make(map[string]*dnsToQuery)
				for{
					select{
					case <-ctx.Done():
					case toQuery <- dnsToQueryChan:
						logger.Debug("will query host ")
						toQueryMap[toQuery.layerDNS.ID] = toQuery
					case <-dnsResChan:
						// { TODO compose up dns response 
						rawIpLayer := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
						rawUdpLayer := packet.Layer(layers.LayerTypeUDP).(*layers.UDP)
						rawDNDLayer := dnsLayer.(*layers.DNS)
						resDNS := rawDNDLayer
						resDNS.QR=true
						resDNS.ID=rawDNDLayer.ID
						resDNS.ANCount=1
						resDNS.OpCode=layers.DNSOpCodeNotify
						resDNS.ResponseCode=layers.DNSResponseCodeNoErr
						resDNS.Answers =[]layers.DNSResourceRecord{
							{
								Class:layers.DNSClassIN,
								Type : layers.DNSTypeA,
								Name:[]byte("www.google.com"),
								Data:[]byte{155,155,155,155},
								IP:net.IPv4(155,155,155,155,),
							},
						}
						ipLayer := &layers.IPv4{
							DstIP:rawIpLayer.SrcIP,
							SrcIP:rawIpLayer.DstIP,
							Version:    4,
							IHL:        5,
							TOS:        0,
							Id:         0,
							Flags:      0,
							FragOffset: 0,
							TTL:        255,
							Protocol:   layers.IPProtocolUDP,
						}
						udpLayer := &layers.UDP{
							SrcPort:rawUdpLayer.DstPort,
							DstPort:rawUdpLayer.SrcPort,
						}
						udpLayer.SetNetworkLayerForChecksum(ipLayer)
						serBuf := gopacket.NewSerializeBuffer()
						err = gopacket.SerializeLayers(serBuf,gopacket.SerializeOptions{FixLengths:true,ComputeCheck
						sums:true},ipLayer,udpLayer,resDNS)
						if err != nil{
							panic(err)
						}
						data = serBuf.Bytes()
						_,err = tunIfce.Write(data)

						// }
						tunIfce.Write(dnsRes)
					}
				}
			}
		)
		asyncCtrl.AsyncDo(
			ctx,
			func(ctx context.Context) {
				var buf = make([]byte, 1024*4)
				for {
					n, err := tun.Read(buf)
					if err != nil {
						logger.Error(err)
						return
					}
					err = tp.WriteIpPacket(buf[:n])
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
		return tunIfce, nil
	}
