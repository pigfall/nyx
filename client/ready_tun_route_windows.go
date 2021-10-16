package client

import (
	"fmt"
	stdnet "net"

	//"fmt"
	"github.com/pigfall/tzzGoUtil/log"
	"github.com/pigfall/tzzGoUtil/net"
	"github.com/pigfall/tzzGoUtil/syscall/winsys"
)

func readyTunRoute(logger log.LoggerLite, serverIp stdnet.IP, tunIfceName string) error {
	err := winsys.LoadIpHelperDLL()
	if err != nil {
		return err
	}
	ifce, err := net.FindIfceByName(tunIfceName)
	if err != nil {
		return err
	}

	return readyTunBase(logger, serverIp, fmt.Sprintf("%d", ifce.Index))
}
