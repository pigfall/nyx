package client

import(
	"github.com/pigfall/tzzGoUtil/net/wintun"
	"fmt"
	"github.com/pigfall/tzzGoUtil/net"

)


func NewTun()(net.TunIfce,error){
	tun,err := wintun.NewTun("yy-ifce",1500)
	if err != nil{
		err = fmt.Errorf("Create tun ifce failed %v",err)
		return nil,err
	}
	return tun,nil
}
