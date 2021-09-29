package proto_handler

import(
	ws "github.com/gorilla/websocket"
	"github.com/pigfall/tzzGoUtil/net"
	"context"
)

func (this *Handler)handleQueryIp(
	ctx context.Context,
	msgBytes []byte,
	conn *ws.Conn,
	unmarshalFunc func([]byte,interface{},
)error){
	ipNet := this.clientVPNIp
}
