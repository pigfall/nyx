package main

import (
	"fmt"
	"context"
	"encoding/json"
	"net"
	stdos "os"
	"flag"
	"github.com/pigfall/yingying/server"
	yy "github.com/pigfall/yingying"
	"github.com/pigfall/tzzGoUtil/log"
	"github.com/pigfall/tzzGoUtil/path/filepath"
	"github.com/pigfall/tzzGoUtil/encoding"
	"github.com/pigfall/tzzGoUtil/os"
)

func main() {
	// { parse cmd line
	var cfgPath string
	flag.StringVar(&cfgPath,"conf","","config file path")
	flag.Parse()
	// }
	if len(cfgPath) == 0 {
		execPath,err := os.GetExecutablePath()
		if err != nil{
			panic(err)
		}
		cfgPath =filepath.Join(execPath,"config.json")
	}

	rawLogger := log.NewJsonLogger(stdos.Stdout)
	loggerMain := log.NewHelper("main",rawLogger,log.LevelDebug)
	loggerMain.Infof("Use config file %s",cfgPath)
	cfg := server.ServeCfg{}
	err := encoding.UnMarshalByFile(cfgPath,&cfg,json.Unmarshal)
	if err != nil{
		loggerMain.Errorf("Parse config file %s failed %v",cfgPath,err)
		stdos.Exit(1)
	}

	loggerMain.Info("Serving")
	ctx := context.Background()
	mode := cfg.Mode
	var tpServer yy.TransportServer
	ipToListen  := net.ParseIP("0.0.0.0")
	if ipToListen ==  nil{
		panic(fmt.Errorf("0.0.0.0 parse failed"))
	}
	switch mode{
	case "udp":
		panic("TODO")
	case  "ws":
		tpServer = server. NewTransportServerWebSocket(ipToListen,cfg.Port)
	default:
		loggerMain.Errorf("Config error, mode is %v undefined , must be udp or ws",mode)
		stdos.Exit(1)
	}
	err = server.Serve(ctx,rawLogger,&cfg,tpServer)
	if err != nil{
		loggerMain.Error(err)
	}
	loggerMain.Info("App quit")
}
