// config package sets global parameters via a configuration file.
//
// On application startup configuration values are read from the file
// ~/.config/pigiron/config.toml and stored to the GlobalParameters
// struct.  An alternate config file may be specified on the command line 
//    pigiron alternate-config.toml
//

package config

import (
	"fmt"
	"strings"
	"strconv"
	"os"
	"path/filepath"
	toml "github.com/pelletier/go-toml"
)

var (
	GlobalParameters = globalParameters {}
	configFilename string
	tomlTree *toml.Tree
)

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
}

// ResetGlobalParameters sets all global configuration parameter to default values."
//
func ResetGlobalParameters() {
	GlobalParameters.OSCServerRoot = "pig"
	GlobalParameters.OSCServerHost = "127.0.0.1"
	GlobalParameters.OSCServerPort = 8000
	GlobalParameters.OSCClientRoot = "pig-client"
	GlobalParameters.OSCClientHost = "127.0.0.1"
	GlobalParameters.OSCClientPort = 8001
	GlobalParameters.OSCClientFilename = ""
	GlobalParameters.MaxTreeDepth = 12
	GlobalParameters.MIDIInputBufferSize = 1024
	GlobalParameters.MIDIInputPollInterval = 0
	GlobalParameters.MIDIOutputBufferSize = 1024
	GlobalParameters.MIDIOutputLatency = 0
}


// determinConfigFilename returns the configuration filename.
// The default filename is <config>/pigiron/config.toml
// where <config> is the User's configuration directory.
// Typically on ~/.config/ on Linux.
//
// An alternate configuration file may be specified on the
// command line.  
// 
func determineConfigFilename() string {
	if len(os.Args) < 2 {
		base, err := os.UserConfigDir()
		if err == nil {
			home, _ := os.UserHomeDir()
			base = filepath.Join(home, ".config")
		}
		return filepath.Join(base, "pigiron", "config.toml")
	} else {
		filename := os.Args[1]
		if len(filename) > 2 && filename[0:1] == "~/" {
			home, _ := os.UserHomeDir()
			filename = filepath.Join(home, filename[2:])
		}
		return filename
	}
}

func splitPath(path string) []string {
	return strings.Split(path, ".")
}


// hasPath returns true if the config file contains path.
//
func hasPath(path string) bool {
	return tomlTree.HasPath(splitPath(path))
}


// readInt reads an int from the config file.
// Returns fallback if path does not exists or its value is invalid.
//
func readInt(path string, fallback int64) int64 {
	if hasPath(path) {
		raw := fmt.Sprintf("%v", tomlTree.Get(path))
		value, err := strconv.Atoi(raw)
		if err != nil {
			msg := "ERROR: Config '%s' expected int, found '%s'. Using default %v\n"
			fmt.Printf(msg, path, raw, fallback)
			return fallback
		} else {
			return int64(value)
		}
	} else {
		msg :="ERROR: Config path '%s' missing, using default %v\n"
		fmt.Printf(msg, path, fallback)
		return fallback
	}
}


// readFloat reads float from config file.
// Returns fallback if path does not exists or its value is invalid.
//
func readFloat(path string, fallback float64) float64 {
	if hasPath(path) {
		raw := fmt.Sprintf("%v", tomlTree.Get(path))
		value, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			msg := "ERROR: Config '%s' expected float, found '%s'. Using default %v\n"
			fmt.Printf(msg, path, raw, fallback)
			return fallback
		} else {
			return value
		}
	} else {
		msg :="ERROR: Config path '%s' missing, using default %v\n"
		fmt.Printf(msg, path, fallback)
		return fallback
	}
}

// readString reads a string value from config file.
// Returns fallback if the path dose not exists.
//
func readString(path string, fallback string) string {
	if hasPath(path) {
		return fmt.Sprintf("%v", tomlTree.Get(path))
	} else {
		msg :="ERROR: Config path '%s' missing, using default %v\n"
		fmt.Printf(msg, path, fallback)
		return fallback
	}
}


func init() {
	ResetGlobalParameters()
	var err error
	configFilename = determineConfigFilename()
	tomlTree, err = toml.LoadFile(configFilename)
	if err != nil {
		fmt.Printf("ERROR: Can not load configuration file \"%s\"\n", configFilename)
		fmt.Println("ERROR:", err.Error())
		fmt.Println()
		return
	} else {
		fmt.Printf("Using config file: \"%s\"\n", configFilename)
		GlobalParameters.OSCServerRoot = readString("osc-server.root", "pig")
		GlobalParameters.OSCServerHost = readString("osc-server.host", "127.0.0.1")
		GlobalParameters.OSCServerPort = readInt("osc-server.port", 8000)
		GlobalParameters.OSCClientRoot = readString("osc-client.root", "pig-client")
		GlobalParameters.OSCClientHost = readString("osc-client.host", "127.0.0.1")
		GlobalParameters.OSCClientPort = readInt("osc-client.port", 8001)
		GlobalParameters.OSCClientFilename = readString("osc-client.file", "")
		GlobalParameters.MaxTreeDepth = readInt("tree.max-depth", 12)
		GlobalParameters.MIDIInputBufferSize = readInt("midi-input.buffer-size", 1024)
		GlobalParameters.MIDIInputPollInterval = readInt("midi-input.poll-interval", 0)
		GlobalParameters.MIDIOutputBufferSize = readInt("midi-output.buffer-size", 1024)
		GlobalParameters.MIDIOutputLatency = readInt("midi-output.latency", 0)
	}
}


// DumpGlobalParameters prints the global configuration values.
//
func DumpGlobalParameters() {
	fmt.Println("Global Parameters:")
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
	