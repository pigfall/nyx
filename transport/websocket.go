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

func(this *tpWebsocket) Read()(msgType yy.TransprtMsgType,data []byte,err error){
	msgTypeInt,data,err := this.conn.ReadMessage()
	if err != nil{
		return 0,nil,err
	}
	if msgTypeInt == ws.BinaryMessage{
		return yy.IpPacket,data,nil
	}
	return yy.Proto,data,nil
}

func(this *tpWebsocket)WriteIpPacket(ipPacketBytes []byte)(error){
	return this.conn.WriteMessage(ws.BinaryMessage,ipPacketBytes)

}

func (this *tpWebsocket) WriteJSON(msg interface{})error{
	return this.conn.WriteJSON(msg)
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

