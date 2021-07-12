package main

import (
	"fmt"
	"os"
	"path/filepath"
	//"time"
	"github.com/plewto/pigiron/config"
	"github.com/plewto/pigiron/osc"
	"github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/op"
	_ "github.com/plewto/pigiron/smf"
)


var banner = []string{
	"    ____  __________________  ____  _   __    ",
	"   / __ \\/  _/ ____/  _/ __ \\/ __ \\/ | / / ",
	"  / /_/ // // / __ / // /_/ / / / /  |/ /     ",
	" / ____// // /_/ // // _, _/ /_/ / /|  /      ",
	"/_/   /___/\\____/___/_/ |_|\\____/_/ |_/     "} 



func printBanner() {
	for _, line := range banner {
		fmt.Println(line)
	}
	fmt.Println()
	fmt.Printf("Version %s\n", config.Version)
	cfig, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("WARNING: Can not dertermin user's config directory.\n")
		fmt.Printf("%s\n", err)
	} else {
		cfig = filepath.Join(cfig, "pigiron")
		fmt.Printf("Configuration directory is '%s'\n", cfig)
	}
	fmt.Println()
}


func main() {
	fmt.Print(config.GlobalParameters.TextColor)
	printBanner()

	

	
	//config.DumpGlobalParameters()
	osc.Init()
	op.Init()
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
