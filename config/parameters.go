package config

import (
	"fmt"
)


/*
 * Defines global configuration parameters.
*/ 

// globalParameters struct holds application wide configuration values.
//
type globalParameters struct {
	OSCServerRoot string
	OSCServerHost string
	OSCServerPort int64
	OSCClientRoot string
	OSCClientHost string
	OSCClientPort int64
	OSCClientFilename string
	MaxTreeDepth int64
	MIDIInputBufferSize int64
	MIDIInputPollInterval int64 // ms
	MIDIOutputBufferSize int64
	MIDIOutputLatency int64
	BannerColor string
	TextColor string
	ErrorColor string
}

// ResetGlobalParameters sets all global configuration parameter to default values."
//
func ResetGlobalParameters() {
	GlobalParameters.OSCServerRoot = "pig"
	GlobalParameters.OSCServerHost = "127.0.0.1"
	GlobalParameters.OSCServerPort = 8020
	GlobalParameters.OSCClientRoot = "pig-client"
	GlobalParameters.OSCClientHost = "127.0.0.1"
	GlobalParameters.OSCClientPort = 8021
	GlobalParameters.OSCClientFilename = ""
	GlobalParameters.MaxTreeDepth = 12
	GlobalParameters.MIDIInputBufferSize = 1024
	GlobalParameters.MIDIInputPollInterval = 0
	GlobalParameters.MIDIOutputBufferSize = 1024
	GlobalParameters.MIDIOutputLatency = 0
	GlobalParameters.BannerColor = getColor("white")
	GlobalParameters.TextColor = getColor("white")
	GlobalParameters.ErrorColor = getColor("red")
}




// DumpGlobalParameters prints the global configuration values.
//
func PrintConfig() {
	fmt.Println("Global configuration values\n")
	fmt.Printf("\tpigiron version: %s\n", Version)
	fmt.Printf("\tconfig file was \"%s\"\n", configFilename)
	fmt.Printf("\tOSCServerRoot         : %v\n", GlobalParameters.OSCServerRoot)
	fmt.Printf("\tOSCServerHost         : %v\n", GlobalParameters.OSCServerHost)
	fmt.Printf("\tOSCServerPort         : %v\n", GlobalParameters.OSCServerPort)
	fmt.Printf("\tOSCClientRoot         : %v\n", GlobalParameters.OSCClientRoot)
	fmt.Printf("\tOSCClientHost         : %v\n", GlobalParameters.OSCClientHost)
	fmt.Printf("\tOSCClientPort         : %v\n", GlobalParameters.OSCClientPort)
	fmt.Printf("\tOSCClientFilename     : %v\n", GlobalParameters.OSCClientFilename)
	fmt.Printf("\tMaxTreeDepth          : %v\n", GlobalParameters.MaxTreeDepth)
	fmt.Printf("\tMIDIInputBufferSize   : %v\n", GlobalParameters.MIDIInputBufferSize)
	fmt.Printf("\tMIDIInputPollInterval : %v\n", GlobalParameters.MIDIInputPollInterval)
	fmt.Printf("\tMIDIOutputBufferSize  : %v\n", GlobalParameters.MIDIOutputBufferSize)
	fmt.Printf("\tMIDIOutputLatency     : %v\n", GlobalParameters.MIDIOutputLatency)
}
