package osc

import (
	//"fmt"
	// goosc "github.com/hypebeast/go-osc/osc"
	//"github.com/rakyll/portmidi"
	//"github.com/plewto/pigiron/midi"
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


// Must execute after config init()
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


func Listen() {
	GlobalServer.ListenAndServe()
}


func Cleanup() {
	GlobalServer.Close()
}













