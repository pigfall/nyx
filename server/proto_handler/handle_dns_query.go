package proto_handler

import(
	"context"
	"time"
"github.com/pigfall/tzzGoUtil/net"
"github.com/pigfall/yingying/proto"
)


func (this *Handler)handleDNSQuery(
	ctx context.Context,
	msgBytes []byte,
	unmarshalFunc func([]byte,interface{})(error),
)(*proto.Msg,interface{}){
	dnsQuery := &proto.DNSQuery{}
	msgRet := &proto.Msg{
		Id:proto.ID_S2C_DNS_QUERY,
	}

	err := unmarshalFunc(msgBytes,dnsQuery)
	if err != nil{
		msgRet.ErrReason = err.Error()
		return msgRet,nil
	}
	// TODO ConCurrent query
	dnsQueryRes := &proto.DNSQueryRes{ID:dnsQuery.ID}
	dnsQueryRes.Answers = make([]*proto.DNSQueryAnswers,0,len(dnsQuery.HostNames))
	for _,hostname := range dnsQuery.HostNames {
		ips,err := net.LookupHost(ctx,time.Second*3,hostname)
		if err != nil {
			msgRet.ErrReason = err.Error()
			msgRet.ErrMsg= err.Error()
			return msgRet,nil
		}
		dnsQueryRes.Answers = append(dnsQueryRes.Answers,&proto.DNSQueryAnswers{
			HostName:hostname,
			Ips:ips,
		})
	}
	return msgRet,dnsQueryRes
}
