package transport

import(
	"encoding/json"
ws	"github.com/gorilla/websocket"

  yy "github.com/pigfall/yingying"
   "github.com/pigfall/yingying/proto"
)


type tpWebsocket struct{
	conn *ws.Conn

}


func NewTPWebSocket(conn *ws.Conn) yy.Transport{
	return &tpWebsocket{
		conn:conn,
	}
} 

func(this *tpWebsocket)WriteIpPacket(ipPacketBytes []byte)(error){
	return this.conn.WriteMessage(ws.BinaryMessage,ipPacketBytes)

}
func(this *tpWebsocket)WriteMsg(msg *proto.Msg,body interface{})(error){
	if body != nil{
		bodyBytes,err := json.Marshal(body)
		if err != nil{
			return err
		}
		msg.Body = bodyBytes
	}
	return this.conn.WriteJSON(msg)
}
