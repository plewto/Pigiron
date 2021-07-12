package config

/*
 * Parse command line arguments.
 * Load toml config file.
 * Establish initial OSC batch filename.
*/

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	toml "github.com/pelletier/go-toml"
	"strings"
	"strconv"
)

const Version = "0.0.1"
	

var (
	GlobalParameters = globalParameters{}
	configFilename string
	BatchFilename string
	tomlTree *toml.Tree
)

// subUserHome substitutes leading '~' character for user home directory.
//
func subUserHome(filename string) string {
	result := filename
	if len(filename) > 0 && string(filename[0]) == "~" {
		home, _ := os.UserHomeDir()
		result = filepath.Join(home, filename[1:])
	}
	return result
}


// parseCommandLine deciphers command line arguments.
//
// --config filename
//     Use alternate configuration file.
//     Defaults to ~/<config-dir>/pigiron/config.toml
//
// --batch filename
//     Sets OSC batch file to run at startup.
//     Defaults to no file.
//
func parseCommandLine() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = ".config"
	}
	// config filename
	defaultFile := filepath.Join(configDir, "pigiron", "config.toml")
	flag.StringVar(&configFilename, "config", defaultFile, "Sets configuration file.")
	configFilename = subUserHome(configFilename)
	// batch filename
	defaultFile = ""
	flag.StringVar(&BatchFilename, "batch", defaultFile, "Sets initial OSC batch file.")
	BatchFilename = subUserHome(BatchFilename)
	
	flag.Parse()
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
// Returns fallback if path does not exists or it's value is invalid.
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
// Returns fallback if path does not exists or it's value is invalid.
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


// readConfigurationFile sets GlobalParameters fields from toml config file.
//
func readConfigurationFile(filename string) {
	var err error
	tomlTree, err = toml.LoadFile(filename)
	if err != nil {
		fmt.Printf("ERROR: Can not open configuration file: '%s'\n", filename)
		fmt.Println("ERROR: ", err.Error())
		fmt.Println()
		ResetGlobalParameters()
		return
	} else {
		fmt.Printf("Using configuration file: '%s'\n", filename)
		GlobalParameters.OSCServerRoot = readString("osc-server.root", "pig")
		GlobalParameters.OSCServerHost = readString("osc-server.host", "127.0.0.1")
		GlobalParameters.OSCServerPort = readInt("osc-server.port", 8020)
		GlobalParameters.OSCClientRoot = readString("osc-client.root", "pig-client")
		GlobalParameters.OSCClientHost = readString("osc-client.host", "127.0.0.1")
		GlobalParameters.OSCClientPort = readInt("osc-client.port", 8021)
		GlobalParameters.OSCClientFilename = readString("osc-client.file", "")
		GlobalParameters.MaxTreeDepth = readInt("tree.max-depth", 12)
		GlobalParameters.MIDIInputBufferSize = readInt("midi-input.buffer-size", 1024)
		GlobalParameters.MIDIInputPollInterval = readInt("midi-input.poll-interval", 0)
		GlobalParameters.MIDIOutputBufferSize = readInt("midi-output.buffer-size", 1024)
		GlobalParameters.MIDIOutputLatency = readInt("midi-output.latency", 0)
		GlobalParameters.TextColor = getColor(readString("colors.text", ""))
		GlobalParameters.ErrorColor = getColor(readString("colors.error", ""))
	}
}


// ResourceFilename returns filename relative to the resources directory.
// The resources directory is located at <config>/resources/
// On Linux this location is ~/.config/pigiron/resources/
//
// Returns non-nil error if resources directory can not be determined.
//
// Example:
// ResourceFilename("foo", "bar.txt") --> ~/.config/pigiron/resources/foo/bar.txt
//
func ResourceFilename(elements ...string) (string, error) {
	cfigdir, err := os.UserConfigDir()
	if err != nil {
		msg := "ERROR: Resource filename can not be determined.\n"
		msg += "ERROR: Can not determine configuration directory location.\n"
		msg += fmt.Sprintf("ERROR: %s\n", err)
		err = fmt.Errorf(msg)
		return "", err
	}
	acc := filepath.Join(cfigdir, "pigiron", "resources")
	for _, e := range elements {
		acc = filepath.Join(acc, e)
	}
	return acc, err
}


func init() {
	defineColors()
	parseCommandLine()
	ResetGlobalParameters()
	readConfigurationFile(configFilename)
}
	

