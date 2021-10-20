package proto_handler

import(
	"fmt"
	"context"
	"encoding/json"
	"github.com/pigfall/yingying/proto"
	"github.com/pigfall/tzzGoUtil/log"
	"github.com/pigfall/tzzGoUtil/net"
)


type Handler struct {
	clientVPNIp *net.IpWithMask
}

func NewHandler(clientVPNIp *net.IpWithMask)*Handler{
	return &Handler{
		clientVPNIp:clientVPNIp,
	}
}



func (this *Handler)Handle(ctx context.Context,rawLogger log.Logger_Log,msg *proto.Msg)(
	res *proto.Msg,
	body interface{},
){
	logger := log.NewHelper("protoHandler",rawLogger,log.LevelDebug)
	var unmarshalFunc  func(bytes []byte,obj interface{})error
	unmarshalFunc = json.Unmarshal
	var handler  func(ctx context.Context,msgBytes []byte,unmarshalFunc func([]byte,interface{})error)(*proto.Msg,interface{})
	switch msg.Id {
	case proto.ID_C2S_QUERY_IP:
		handler = this.handleQueryIp
	case proto.ID_C2S_DNS_QUERY:
		handler = this.handleDNSQuery
	default :
		err := fmt.Errorf("Undefined Msg Id %v",msg.Id)
		logger.Error(err)
		return nil,err
	}
	return handler(ctx,msg.Body,unmarshalFunc)
}
