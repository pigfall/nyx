package proto_handler

import(
	"context"
"github.com/pigfall/yingying/proto"
)

func (this *Handler)handleQueryIp(
	ctx context.Context,
	msgBytes []byte,
	unmarshalFunc func([]byte,interface{})(error),
)(*proto.Msg,interface{}){
	ipNet := this.clientVPNIp
	return &proto.Msg{
		Id:proto.ID_S2C_QUERY_IP,
	},&proto.S2C_ClientVPNIpNet{
		IpNet:ipNet.FormatAsIpSlashMask(),
	}
}
