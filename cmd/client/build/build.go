package main

import (
	"path"
	goos "os"
	"log"
	"github.com/pigfall/tzzGoUtil/os"
	"github.com/pigfall/tzzGoUtil/process"
)

func init(){
	log.SetFlags(log.LstdFlags)
}


func main() {
	buildBinPath,err := os.GetExecutablePath()
	if err != nil{
		log.Println("Get build script path failed ",err)
		goos.Exit(1)
	}
	buildBinDir := path.Dir(buildBinPath)
	log.Println("Working directory ", buildBinDir)
	clientBinDir:= path.Join(buildBinDir,"../")
	err = goos.Chdir(clientBinDir)
	if err != nil{
		log.Printf("change to directory %s failed\n",clientBinDir)
	}
	log.Printf("Change to directory %s\n",clientBinDir)
	out,errOut,err := process.ExeOutput("go","build","-o","client.exe",".")
	log.Println("Compiling")
	if err !=nil{
		log.Println("Compiled failed: %v, %v, %v",out,errOut,err)
		goos.Exit(1)
	}
}
