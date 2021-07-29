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
	"github.com/plewto/pigiron/piglog"
	_ "github.com/plewto/pigiron/help"
	
)


var banner = []string{
	"    ____  __________________  ____  _   __    ",
	"   / __ \\/  _/ ____/  _/ __ \\/ __ \\/ | / / ",
	"  / /_/ // // / __ / // /_/ / / / /  |/ /     ",
	" / ____// // /_/ // // _, _/ /_/ / /|  /      ",
	"/_/   /___/\\____/___/_/ |_|\\____/_/ |_/     "} 



func printBanner() {
	fmt.Printf("\n")
	fmt.Print(config.GlobalParameters.BannerColor)
	for _, line := range banner {
		fmt.Println(line)
	}
	fmt.Print("\n")
	cfig, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("WARNING: Can not dertermin user's config directory.\n")
		fmt.Printf("%s\n", err)
	} else {
		cfig = filepath.Join(cfig, "pigiron")
		fmt.Printf("Configuration directory is '%s'\n", cfig)
	}
	fmt.Printf("Configuration file: %s\n", config.ConfigFilename())
	if config.GlobalParameters.EnableLogging {
		fmt.Printf("Logging to: %s\n", piglog.Logfile())
	} else {
		fmt.Printf("Logging disabled\n")
	}
	fmt.Printf("Version: %s\n", config.Version)
	

	fmt.Print(config.GlobalParameters.TextColor)
	fmt.Print("\n\n")
}


func main() {
	piglog.Log("-------- Pigiron main()")
	printBanner()
	piglog.Log(config.ConfigInfo())
	osc.Init()
	op.Init()
	osc.Listen()
	if config.BatchFilename != "" {  
		// osc.LoadBatchFile(config.BatchFilename) // BUG 003 -- do not use --
		fmt.Println("BUG 003, command line batch file is disabled.")
	}
	go osc.REPL()
	fmt.Println()
	// main loop
	for { 
		if osc.Exit {
			Exit()
		}
		op.ProcessInputs()
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
	piglog.Close()
}
