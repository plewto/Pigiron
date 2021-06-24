package main

import (
	"fmt"
	"os"
	//"time"
	"github.com/plewto/pigiron/config"
	"github.com/plewto/pigiron/osc"
	"github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/op"
)


func main() {
	fmt.Print(config.GlobalParameters.TextColor)
	fmt.Println("Welcome to Pigiron")
	config.DumpGlobalParameters()
	osc.Init()
	osc.Listen()
	if config.BatchFilename != "" {  
		// osc.LoadBatchFile(config.BatchFilename) // BUG 003 -- do not use --
		fmt.Println("BUG 003, command line batch file is disabled.")
	}
	go osc.REPL()
	fmt.Println()
	for { // main loop
		if osc.Exit {
			Exit()
		}
	}
}

func Exit() {
	fmt.Println("Pigiron exit.")
	Cleanup()
	os.Exit(0)
}



func Cleanup() {
	midi.Cleanup()
	op.Cleanup()
	osc.Cleanup()
	
}
