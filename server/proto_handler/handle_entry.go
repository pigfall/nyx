package proto_handler

import(
	"fmt"
	"context"
	"encoding/json"
	ws "github.com/gorilla/websocket"
	"github.com/pigfall/nyx/proto"
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



func (this *Handler)Handle(ctx context.Context,rawLogger log.Logger_Log,conn *ws.Conn,msg *proto.Msg)(error){
	logger := log.NewHelper("protoHandler",rawLogger,log.LevelDebug)
	var unmarshalFunc  func(bytes []byte,obj interface{})error
	unmarshalFunc = json.Unmarshal
	var handler  func(ctx context.Context,msgBytes []byte,conn *ws.Conn,unmarshalFunc func([]byte,interface{})error)error
	switch msg.Id {
	case proto.ID_C2S_QUERY_IP:
		handler = this.handleQueryIp
	default:
		err := fmt.Errorf("Undefined Msg Id %v",msg.Id)
		logger.Error(err)
		return err
	}
	return handler(ctx,msg.Body,conn,unmarshalFunc)
}
