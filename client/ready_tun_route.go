package client

import (
	"fmt"
	stdnet "net"

	"github.com/pigfall/tzzGoUtil/log"
	"github.com/pigfall/tzzGoUtil/net"
)

func readyTunBase(logger log.LoggerLite, serverIp stdnet.IP, tunIfceName string) error {
	logger.Info("tunIfceName ", tunIfceName)
	defaultRoute, err := net.GetDefaultRouteRule()
	if err != nil {
		err = fmt.Errorf("Get default route from kernel failed: %w", err)
		logger.Error(err)
		return err
	}
	net.DelRoute(serverIp)
	defaultRouteIfce, err := net.FindIfceByName(defaultRoute.DevName)
	if err != nil {
		logger.Error(err)
		return err
	}
	err = net.AddRoute(serverIp, fmt.Sprintf("%d", defaultRouteIfce.Index), defaultRoute.Via)
	if err != nil {
		logger.Error(err)
		return err
	}
	targetA, err := net.FromIpSlashMask("0.0.0.0/1")
	if err != nil {
		logger.Error(err)
		return err
	}
	targetB, err := net.FromIpSlashMask("128.0.0.0/1")
	if err != nil {
		logger.Error(err)
		return err
	}
	err = net.AddRouteIpNet(targetA, tunIfceName, nil)
	if err != nil {
		logger.Error("Set route table failed %v", err)
		return err
	}
	err = net.AddRouteIpNet(targetB, tunIfceName, nil)
	if err != nil {
		logger.Error("Set route table failed %v", err)
		return err
	}

	return nil
}
