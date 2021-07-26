package config

import "fmt"


/*
 * Defines global configuration parameters.
*/ 

// globalParameters struct holds application wide configuration values.
//
type globalParameters struct {
	EnableLogging bool
	Logfile string
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
	GlobalParameters.EnableLogging = true
	GlobalParameters.Logfile = "!/log"
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


// PrintConfig prints the global configuration values.
//
func PrintConfig() {
	fmt.Println(ConfigInfo())
}

// ConfigInfo returns string representation of global configuration values.
//
func ConfigInfo() string {
	acc := "Global configuration values\n"
	acc += fmt.Sprintf("\tpigiron version: %s\n", Version)
	acc += fmt.Sprintf("\tconfig file was \"%s\"\n", configFilename)
	acc += fmt.Sprintf("\tEnableLogging         : %v\n", GlobalParameters.EnableLogging)
	acc += fmt.Sprintf("\tLogfile               : %v\n", GlobalParameters.Logfile)
	acc += fmt.Sprintf("\tOSCServerRoot         : %v\n", GlobalParameters.OSCServerRoot)
	acc += fmt.Sprintf("\tOSCServerHost         : %v\n", GlobalParameters.OSCServerHost)
	acc += fmt.Sprintf("\tOSCServerPort         : %v\n", GlobalParameters.OSCServerPort)
	acc += fmt.Sprintf("\tOSCClientRoot         : %v\n", GlobalParameters.OSCClientRoot)
	acc += fmt.Sprintf("\tOSCClientHost         : %v\n", GlobalParameters.OSCClientHost)
	acc += fmt.Sprintf("\tOSCClientPort         : %v\n", GlobalParameters.OSCClientPort)
	acc += fmt.Sprintf("\tOSCClientFilename     : %v\n", GlobalParameters.OSCClientFilename)
	acc += fmt.Sprintf("\tMaxTreeDepth          : %v\n", GlobalParameters.MaxTreeDepth)
	acc += fmt.Sprintf("\tMIDIInputBufferSize   : %v\n", GlobalParameters.MIDIInputBufferSize)
	acc += fmt.Sprintf("\tMIDIInputPollInterval : %v\n", GlobalParameters.MIDIInputPollInterval)
	acc += fmt.Sprintf("\tMIDIOutputBufferSize  : %v\n", GlobalParameters.MIDIOutputBufferSize)
	acc += fmt.Sprintf("\tMIDIOutputLatency     : %v\n", GlobalParameters.MIDIOutputLatency)
	return acc
}
