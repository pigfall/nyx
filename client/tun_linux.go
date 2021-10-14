package client

import(
	water_wrap "github.com/pigfall/tzzGoUtil/net/water_tun_wrap"
	"fmt"
	"github.com/pigfall/tzzGoUtil/net"

)


func NewTun()(net.TunIfce,error){
	tun,err := water_wrap.NewTun()
	if err != nil{
		err = fmt.Errorf("Create tun ifce failed %v",err)
		return nil,err
	}
	return tun,nil
}
