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

	// Exit application if true.
	Exit bool = false
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
	
	// Create global OSC server
	host = config.GlobalParameters.OSCServerHost
	port = int(config.GlobalParameters.OSCServerPort)
	root = config.GlobalParameters.OSCServerRoot
	GlobalServer = NewServer(host, port, root)
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

