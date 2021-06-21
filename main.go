package main

import (
	"fmt"
	"os"

	"github.com/plewto/pigiron/config"
	"github.com/plewto/pigiron/osc"
	"github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/op"
)


func main() {
	fmt.Println("Pigiron")
	config.DumpGlobalParameters()
	osc.Init()
	osc.Listen()
	go repl()
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
