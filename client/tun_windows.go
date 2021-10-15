package client

import(
	"github.com/pigfall/tzzGoUtil/net/wintun"
	wg "github.com/pigfall/wtun-go"
	"fmt"
	"github.com/pigfall/tzzGoUtil/net"

)


func NewTun()(net.TunIfce,error){
	tun,err := wg.NewTun("yy-ifce")
	if err != nil{
		return nil,err
	}
	return tun,nil
}
