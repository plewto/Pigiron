package osc

import (
	"fmt"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/config"
)

var (
	globalClient *PigClient
	globalServer *PigServer
	empty []string
	Exit bool = false  // exit application if true
)


func init() {
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
}

type PigServer struct {
	backingServer *goosc.Server
	dispatcher *goosc.StandardDispatcher
	root string
	client *PigClient
}


func NewServer(ip string, port int, root string) *PigServer {
	server := new(PigServer)
	server.root = root
	server.client = globalClient
	addr := fmt.Sprintf("%s:%d", ip, port)
	server.dispatcher = goosc.NewStandardDispatcher()
	server.backingServer = &goosc.Server {
		Addr: addr,
		Dispatcher: server.dispatcher,
	}
	server.AddMsgHandler("ping", server.ping)
	server.AddMsgHandler("exit", server.exit)
	server.AddMsgHandler("new-operator", server.newOperator)
	server.AddMsgHandler("new-midi-input", server.newMIDIInput)
	server.AddMsgHandler("new-midi-output", server.newMIDIOutput)
	server.AddMsgHandler("delete-operator", server.deleteOperator)
	server.AddMsgHandler("connect", server.connect)
	server.AddMsgHandler("disconnect",     server.disconnect)
	server.AddMsgHandler("disconnect-all", server.disconnectAll)
	server.AddMsgHandler("destroy-forest", server.destroyForest)
	server.AddMsgHandler("print-forest", server.printForest)
	server.AddMsgHandler("q-is-parent", server.queryIsParent)
	server.AddMsgHandler("q-midi-inputs", server.queryMIDIInputs)
	server.AddMsgHandler("q-midi-outputs", server.queryMIDIOutputs)
	server.AddMsgHandler("q-operators", server.queryOperators)
	server.AddMsgHandler("q-roots", server.queryRoots)
	server.AddMsgHandler("q-children", server.queryChildren)
	server.AddMsgHandler("q-parents", server.queryParents)
	server.AddMsgHandler("panic", server.panic)
	server.AddMsgHandler("reset", server.panic)
	return server
}

func (s *PigServer) AddMsgHandler(command string, handler func(msg *goosc.Message)) {
	addr := fmt.Sprintf("/%s/%s", s.root, command)
	s.dispatcher.AddMsgHandler(addr, handler)
}

func (s *PigServer) ListenAndServe() {
	fmt.Println("OSC Listening....")
	go s.backingServer.ListenAndServe()
}

func Listen() {
	globalServer.ListenAndServe()
}




