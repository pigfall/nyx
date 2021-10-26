package main

import(
	"encoding/json"
	"github.com/pigfall/tzzGoUtil/path/filepath"
	osHelper "github.com/pigfall/tzzGoUtil/os"
	"github.com/pigfall/tzzGoUtil/encoding"
	"flag"
	"os"
 "github.com/pigfall/yingying/client"
"github.com/pigfall/tzzGoUtil/log"
)


func main() {
	// {}
	var cfgPath string
	flag.StringVar(&cfgPath,"conf","","config file path")
	flag.Parse()
	// }
	if len(cfgPath) == 0 {
		execPath,err := osHelper.GetExecutablePath()
		if err != nil{
			panic(err)
		}
		cfgPath =filepath.Join(execPath,"config.json")
	}

	//
	// ctx:= context.Background()
	rawLogger := log.NewJsonLogger(os.Stdout)
	loggerMain := log.NewHelper("main",rawLogger,log.LevelDebug)
	cfg := client.RunCfg{}
	err := encoding.UnMarshalByFile(cfgPath,&cfg,json.Unmarshal)
	if err != nil{
		loggerMain.Errorf("Parse config file %s failed %v",cfgPath,err)
		mainQuit(loggerMain)
	}
	err = client.LoopRun(&cfg)
	// err = client.Run(ctx,rawLogger,&cfg)
	if err != nil{
		loggerMain.Error(err)
		mainQuit(loggerMain)
	}
	loggerMain.Info("App quit")
}
