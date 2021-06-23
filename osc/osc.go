package osc

import (
	"fmt"

	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/rakyll/portmidi"
	"github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/util"
	"github.com/plewto/pigiron/config"
)

var (
	globalClient PigClient
	globalServer PigServer
	empty []string

	// Exit application if true.
	// This should proabbly be replaced with a go channel message.
	Exit bool = false
)

const (
	XpString util.ExpectType = util.XpString
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

	AddOSCHandler(globalServer, "ping", remotePing)
	AddOSCHandler(globalServer, "exit", remoteExit)
	AddOSCHandler(globalServer, "q-midi-inputs", remoteQueryMIDIInputs)
	AddOSCHandler(globalServer, "q-midi-outputs", remoteQueryMIDIOutputs)
	
}


func Listen() {
	globalServer.ListenAndServe()
}


func Cleanup() {
	globalServer.Close()
}


// osc /pig/ping -> ACK
// diagnostic function.
//
func remotePing(msg *goosc.Message)([]string, error) {
	var err error
	fmt.Printf("PING %s\n", msg.Address)
	return empty, err
}


// osc /pig/exit -> ACK
// Terminate application
//
func remoteExit(msg *goosc.Message)([]string, error) {
	var err error
	Exit = true
	return empty, err
}


// osc /pig/q-midi-inputs
// -> ACK list of MIDI input devices
//
func remoteQueryMIDIInputs(msg *goosc.Message)([]string, error) {
	var err error
	ids := midi.InputIDs()
	acc := make([]string, len(ids))
	fmt.Println("MIDI Input devices:")
	for i, id := range ids {
		info := portmidi.Info(id)
		fmt.Printf("\t%s\n", info.Name)
		acc[i] = fmt.Sprintf("\"%s\" ", info.Name)
	}
	return acc, err
}


// osc /pig/q-midi-outputs
// -> ACK list of MIDI output devices
//
func remoteQueryMIDIOutputs(msg *goosc.Message)([] string, error) {
	var err error
	ids := midi.OutputIDs()
	acc := make([]string, len(ids))
	fmt.Println("MIDI Output devices:")
	for i, id := range ids {
		info := portmidi.Info(id)
		fmt.Printf("\t%s\n", info.Name)
		acc[i] = fmt.Sprintf("\"%s\" ", info.Name)
	}
	return acc, err
}
