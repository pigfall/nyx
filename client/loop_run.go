package client

import(
	"errors"
	"github.com/pigfall/tzzGoUtil/log"
	"os"
	"context"
		golog "log"
)


func LoopRun(cfg *RunCfg) error{
	golog.Println("Loop Running")
	rawLogger := log.NewJsonLogger(os.Stdout)
	var err error
	for {
		ctx,cancel := context.WithCancel(context.Background())
		err = Run(ctx,rawLogger,cfg)
		golog.Println("RunOnce quit")
		cancel()
		var errQuitIns *errQuit 
		if errors.As(err,&errQuitIns){
			break
		}
		golog.Println("Retry")
	}
	golog.Println("Client quit ",err)
	return err
}
