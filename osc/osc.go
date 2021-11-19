package osc

import (
	"github.com/plewto/pigiron/config"
	"github.com/plewto/pigiron/pigpath"
)

var (
	globalResponder Responder
	replResponder Responder
	GlobalServer PigServer
	empty []string
	commands map[string]bool
	Exit bool = false  // Flag to main-loop, if true exit application.
)


// Init() initializes the osc package.
// Init() must not be called prior to initialization of the config package.
//
func Init() {
	// Create global responders
	host := config.GlobalParameters.OSCClientHost
	port := int(config.GlobalParameters.OSCClientPort)
	root := config.GlobalParameters.OSCClientRoot
	filename := pigpath.SubSpecialDirectories(config.GlobalParameters.OSCClientFilename)
	globalResponder = NewBasicResponder(host, port, root, filename)
	replResponder = NewREPLResponder()
	commands = make(map[string]bool)
	// Create global OSC server
	host = config.GlobalParameters.OSCServerHost
	port = int(config.GlobalParameters.OSCServerPort)
	root = config.GlobalParameters.OSCServerRoot
	GlobalServer = NewServer(host, port, root)
	AddHandler(GlobalServer, "exec", remoteEval)
}


func isCommand(s string) bool {
	_, flag := commands[s]
	return flag || len(s) == 0
}

// Listen() starts OSC server.
//
func Listen() {
	GlobalServer.ListenAndServe()
}

// Cleanup() closes OSC server.
// Cleanup should only be called on application termination.
//
func Cleanup() {
	GlobalServer.Close()
}

