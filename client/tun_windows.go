package client

import(
	wg "github.com/pigfall/wtun-go"
	"github.com/pigfall/tzzGoUtil/net"

)

type Tun struct{
	*wg.Tun
}

func NewTun()(net.TunIfce,error){
	err :=wg.InitWinTun("wintun.dll")
	if err != nil{
		return nil,err
	}
	tun,err := wg.NewTun("yy-ifce")
	if err != nil{
		return nil,err
	}
	return &Tun{Tun:tun},nil
}

func (this *Tun) SetIp(ip ...string)error{
	ipNet,err := net.FromIpSlashMask(ip[0])
	if err != nil{
		return err
	}
	return this.Tun.SetIp(ipNet)
}

