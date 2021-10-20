package client

import (
	"context"
	"fmt"
	stdnet "net"

	"github.com/pigfall/tzzGoUtil/async"
	"github.com/pigfall/yingying/proto"
	"github.com/pigfall/tzzGoUtil/log"
	"github.com/pigfall/tzzGoUtil/net"
	yy "github.com/pigfall/yingying"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket"
	// water_wrap "github.com/pigfall/tzzGoUtil/net/water_tun_wrap"
)

func readyTun(
	ctx context.Context,
	logger log.LoggerLite,
	tunIp *net.IpWithMask,
	asyncCtrl *async.Ctrl,
	tp yy.Transport,
	serverIp stdnet.IP,
	dnsQueryResChan chan *proto.DNSQueryRes,
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
	var dnsToQueryChan =make(chan *dnsToQuery,10)
	asyncCtrl.AsyncDo(
			ctx,
			func(ctx context.Context){
				var toQueryMap  =make(map[string]*dnsToQuery)
				OUT:
				for{
					select{
					case <-ctx.Done():
					case toQuery :=<- dnsToQueryChan:
						logger.Debug("will query host ")
						layerDNS := toQuery.layerDNS
						toQueryMap[fmt.Sprintf("%d",toQuery.layerDNS.ID)] = toQuery
						// TODO
						if len(layerDNS.Questions) == 0{
							logger.Error("DNS Questions is nil, not to query from server")
							continue
						}
						hostNames := make([]string,0,len(layerDNS.Questions))
						for _,question := range layerDNS.Questions{
							hostNames = append(hostNames,string(question.Name))
						}
						err = tp.WriteMsg(
							&proto.Msg{
								Id:proto.ID_C2S_DNS_QUERY,
							},
							&proto.DNSQuery{
								ID:fmt.Sprintf("%d",toQuery.layerDNS.ID),
								HostNames:hostNames,
							},
						)
						if err != nil {
							logger.Error("Failed to send dnsQuery to server", err)
							continue
						}
					case dnsQueryRes := <-dnsQueryResChan:
						// { TODO compose up dns response 
						dnsQuery := toQueryMap[dnsQueryRes.ID]
						if dnsQuery == nil{
							logger.Error("Not found dns id ", dnsQueryRes.ID)
							continue
						}
						packet := dnsQuery.packet
						dnsLayer := dnsQuery.layerDNS
						rawIpLayer := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
						rawUdpLayer := packet.Layer(layers.LayerTypeUDP).(*layers.UDP)
						rawDNDLayer := dnsLayer
						resDNS := rawDNDLayer
						resDNS.QR=true
						resDNS.ID=rawDNDLayer.ID
						resDNS.ANCount = uint16(len(dnsQueryRes.Answers))
						// resDNS.ANCount=1
						resDNS.OpCode=layers.DNSOpCodeNotify
						resDNS.ResponseCode=layers.DNSResponseCodeNoErr
						resDNS.Answers = make([]layers.DNSResourceRecord,0,resDNS.ANCount)
						for _,answerFromServer := range dnsQueryRes.Answers {
							hostIP := stdnet.ParseIP(answerFromServer.Ips[0])
							if hostIP == nil{
								logger.Errorf("invalid ip format %s ",answerFromServer.Ips[0])
								continue OUT
							}
							answer := &layers.DNSResourceRecord{
								Class:layers.DNSClassIN,
								Name:[]byte(answerFromServer.HostName),
								Data : hostIP,
								IP : hostIP,
							}
							if net.IsIpv4(hostIP){
								answer.Type = layers.DNSTypeA
							}else{
								answer.Type = layers.DNSTypeAAAA
							}
							resDNS.Answers = append(resDNS.Answers,*answer)
						}
						//resDNS.Answers =[]layers.DNSResourceRecord{
						//	{
						//		Class:layers.DNSClassIN,
						//		Type : layers.DNSTypeA,
						//		Name:[]byte("www.google.com"),
						//		Data:[]byte{155,155,155,155},
						//		IP:net.IPv4(155,155,155,155,),
						//	},
						//}
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
						err = gopacket.SerializeLayers(serBuf,gopacket.SerializeOptions{FixLengths:true,ComputeChecksums:true},ipLayer,udpLayer,resDNS)
						if err != nil{
							panic(err)
						}
						data := serBuf.Bytes()
						_,err = tunIfce.Write(data)
						// }
					}
				}
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
					if dnsQuery:=isDNSReq(data,logger);dnsQuery != nil{
						dnsToQueryChan <- dnsQuery
						continue
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
		return tunIfce, nil
	}
	
	func isDNSReq(ipPacket []byte,logger log.LoggerLite)(*dnsToQuery){
		packet := gopacket.NewPacket(ipPacket,layers.LayerTypeIPv4,gopacket.Default)
		if errLayer:= packet.ErrorLayer();errLayer != nil{
			logger.Error("parse ipv4 packet failed: ",errLayer)
			return nil
		}
		dnsLayer := packet.Layer(layers.LayerTypeDNS)
		if dnsLayer != nil{
			return &dnsToQuery{
				packet:packet,
				layerDNS:(dnsLayer).(*layers.DNS),
			}
		}

		return nil
	}
