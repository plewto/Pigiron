package osc

import (
	"fmt"
	goosc "github.com/hypebeast/go-osc/osc"
	
)

type PigServer interface {
	Root() string
	SetRoot(string)
	Client() PigClient
	SetClient(PigClient)
	AddMsgHandler(command string, handler func(msg *goosc.Message))
	ListenAndServe()
	IP() string
	Port() int
	//Close()
}

type OSCServer struct {
	backingServer *goosc.Server
	dispatcher *goosc.StandardDispatcher
	root string
	client PigClient
	ip string
	port int
}


func NewServer(ip string, port int, root string) PigServer {
	server := new(OSCServer)
	server.ip = ip
	server.port = port
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
	server.AddMsgHandler("help", server.help)
	return server
}

func (s *OSCServer) AddMsgHandler(command string, handler func(msg *goosc.Message)) {
	addr := fmt.Sprintf("/%s/%s", s.root, command)
	s.dispatcher.AddMsgHandler(addr, handler)
}

func (s *OSCServer) ListenAndServe() {
	fmt.Printf("OSC Listening: %s:%d  /%s\n", s.ip, s.port, s.root)
	go s.backingServer.ListenAndServe()
}

func (s *OSCServer) Root() string {
	return s.root
}

func (s *OSCServer) SetRoot(root string) {
	s.root = root
}

func (s *OSCServer) Client() PigClient {
	return s.client
}

func (s *OSCServer) SetClient(client PigClient) {
	s.client = client
}

func (s *OSCServer) IP() string {
	return s.ip
}

func (s *OSCServer) Port() int {
	return s.port
}

// func (s *OSCServer) Close() {
// 	s.backingServer.Close()
// }
