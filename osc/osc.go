package osc

import (
	"fmt"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/rakyll/portmidi"
	"github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/config"
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
	filename := config.GlobalParameters.OSCClientFilename
	globalResponder = NewBasicResponder(host, port, root, filename)
	replResponder = NewREPLResponder()
	
	// Create global OSC server
	host = config.GlobalParameters.OSCServerHost
	port = int(config.GlobalParameters.OSCServerPort)
	root = config.GlobalParameters.OSCServerRoot
	GlobalServer = NewServer(host, port, root)
	AddHandler(GlobalServer, "ping", remotePing)
	AddHandler(GlobalServer, "exit", remoteExit)
	AddHandler(GlobalServer, "q-midi-inputs", remoteQueryMIDIInputs)
	AddHandler(GlobalServer, "q-midi-outputs", remoteQueryMIDIOutputs)
	AddHandler(GlobalServer, "batch", remoteBatchLoad)
	AddHandler(GlobalServer, "q-commands", remoteQueryCommands)
}


func Listen() {
	GlobalServer.ListenAndServe()
}


func Cleanup() {
	GlobalServer.Close()
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
	for i, id := range ids {
		info := portmidi.Info(id)
		acc[i] = fmt.Sprintf("\"%s\" ", info.Name)
	}
	return acc, err
}


// osc /pig/q-midi-outputs
// -> ACK list of MIDI output devices
//
func remoteQueryMIDIOutputs(msg *goosc.Message)([]string, error) {
	var err error
	ids := midi.OutputIDs()
	acc := make([]string, len(ids))
	for i, id := range ids {
		info := portmidi.Info(id)
		acc[i] = fmt.Sprintf("\"%s\" ", info.Name)
	}
	return acc, err
}

// osc /pig/batch filename
// 
//
func remoteBatchLoad(msg *goosc.Message)([]string, error) {
	template := []ExpectType{XpString}
	args, err := Expect(template, msg.Arguments)
	if err != nil {
		fmt.Print(config.GlobalParameters.ErrorColor)
		fmt.Printf("ERROR: %s\n", msg.Address)
		fmt.Printf("ERROR: %s\n", err)
		fmt.Print(config.GlobalParameters.TextColor)
		return empty, err
	}
	filename := fmt.Sprintf("%s", args[0])
	err = BatchLoad(filename)
	return empty, err
}


func remoteQueryCommands(msg *goosc.Message)([]string, error) {
	var err error
	return GlobalServer.Commands(), err
}
