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
	Close()
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
	//++ server.AddMsgHandler("new-operator", server.remoteNewOperator)
	//++ server.AddMsgHandler("new-midi-input", server.remoteNewMIDIInput)
	//++ server.AddMsgHandler("new-midi-output", server.remoteNewMIDIOutput)
	//++ server.AddMsgHandler("delete-operator", server.remoteDeleteOperator)
	//++ server.AddMsgHandler("destroy-forest", server.remoteDestroyForest)
	//++ server.AddMsgHandler("print-forest", server.remotePrintForest)
	//++ server.AddMsgHandler("q-operators", server.remoteQueryOperators)
	//++ server.AddMsgHandler("q-roots", server.remoteQueryRoots)
	//++ server.AddMsgHandler("panic", server.remotePanic)
	//++ server.AddMsgHandler("reset", server.remotePanic)
	//server.AddMsgHandler("connect", server.remoteConnect)
	//server.AddMsgHandler("disconnect",     server.remoteDisconnect)
	//server.AddMsgHandler("disconnect-all", server.remoteDisconnectAll)
	//server.AddMsgHandler("q-is-parent", server.remoteQueryIsParent)
	//server.AddMsgHandler("q-children", server.remoteQueryChildren)
	//server.AddMsgHandler("q-parents", server.remoteQueryParents)
	//server.AddMsgHandler("help", server.remoteHelp)
	//server.AddMsgHandler("batch", server.remoteBatchLoad)

	

	


	//++ server.AddMsgHandler("q-midi-inputs", server.remoteQueryMIDIInputs)
	//++ server.AddMsgHandler("q-midi-outputs", server.remoteQueryMIDIOutputs)

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

func (s *OSCServer) Close() {
	s.backingServer.CloseConnection()
}


func AddOSCHandler(s PigServer, address string, handler func(*goosc.Message)([]string, error)) {

	address = fmt.Sprintf("/%s/%s", s.Root(), address)
	fmt.Printf("DEBUG address --> %s\n", address)
	
	var result = func(msg *goosc.Message) {
		status, err := handler(msg)
		if err != nil {
			s.Client().Error(address, status)
		} else {
			s.Client().Ack(address, status)
		}
	}
	s.AddMsgHandler(address, result)
}
