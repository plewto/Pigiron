package osc

import (
	"fmt"
	"github.com/plewto/pigiron/config"
)

var (
	globalClient PigClient
	globalServer PigServer
	replServer PigServer
	empty []string

	// Exit application if true.
	// This should proabbly be replaced with a go channel message.
	Exit bool = false
)




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
	// Create REPL server
	host = config.GlobalParameters.REPLHost
	port = int(config.GlobalParameters.REPLPort)
	root = config.GlobalParameters.REPLRoot
	replServer = NewServer(host, port, root)
	replServer.SetClient(REPLClient{})
}


func prompt() {
	root := config.GlobalParameters.OSCServerRoot
	fmt.Printf("\n/%s/ ", root)
}


func Listen() {
	globalServer.ListenAndServe()
	replServer.ListenAndServe()
}


func Cleanup() {
	//globalServer.Close()
}

