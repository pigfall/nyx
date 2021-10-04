package main

import (
	"context"
	"encoding/json"
	stdos "os"
	"flag"
	"github.com/pigfall/yingying/server"
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
	err = server.Serve(ctx,rawLogger,&cfg)
	if err != nil{
		loggerMain.Error(err)
	}
}
