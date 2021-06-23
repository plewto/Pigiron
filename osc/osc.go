package osc

import (
	"github.com/plewto/pigiron/config"
)

var (
	globalClient PigClient
	globalServer PigServer
	empty []string

	// Exit application if true.
	// This should proabbly be replaced with a go channel message.
	Exit bool = false

	errorFlag bool = false
)


func ClearError() {
	errorFlag = false
}

func signalError() {
	errorFlag = true
}

func OSCError() bool {
	return errorFlag
}


// Must execute after config init()
//
func Init() {
	// Create global OSC client
	host := config.GlobalParameters.OSCClientHost
	port := int(config.GlobalParameters.OSCClientPort)
	root := config.GlobalParameters.OSCClientRoot
	filename := config.GlobalParameters.OSCClientFilename
	globalClient = NewClient(host, port, root, filename)
	// Create global OSC server
	host = config.GlobalParameters.OSCServerHost
	port = int(config.GlobalParameters.OSCServerPort)
	root = config.GlobalParameters.OSCServerRoot
	globalServer = NewServer(host, port, root)

	AddOSCHandler(globalServer, "ping", remotePing)
	
}



func Listen() {
	globalServer.ListenAndServe()
}


func Cleanup() {
	globalServer.Close()
}


