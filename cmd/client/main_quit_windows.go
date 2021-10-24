package main

import(
	"os"
	"github.com/pigfall/tzzGoUtil/log"
	"fmt"
	"bufio"
)


func mainQuit(logger log.LoggerLite){
	logger.Info("App quit")
	rd := bufio.NewReader(os.Stdin)
	fmt.Println("Press enter to close from console")
	_,_ = rd.ReadString('\n')
}
