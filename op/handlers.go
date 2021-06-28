package op

import (
	"fmt"
	//"strings"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/osc"
)

var (
	empty []string
)

// Add general op-related handlers to global OSC server
//
func Init() {
	fmt.Println("opserver.Init() executing -- REMOVE THIS LIEN")
	server := osc.GlobalServer
	osc.AddHandler(server, "new-operator", remoteNewOperator)
	osc.AddHandler(server, "new-midi-input", remoteNewMIDIInput)
	osc.AddHandler(server, "new-midi-output", remoteNewMIDIOutput)
	osc.AddHandler(server, "delete-operator", remoteDeleteOperator)
	osc.AddHandler(server, "delete-all-operators", remoteDeleteAllOperators)
	osc.AddHandler(server, "connect", remoteConnect)
	osc.AddHandler(server, "disconnect-child", remoteDisconnect)
	osc.AddHandler(server, "disconnect-all", remoteDisconnectAll)
	osc.AddHandler(server, "disconnect-parents", remoteDisconnectParents)
        osc.AddHandler(server, "q-operator-types", remoteQueryOperatorTypes)
	osc.AddHandler(server, "q-operators", remoteQueryOperators)
	osc.AddHandler(server, "q-roots", remoteQueryRoots)
	osc.AddHandler(server, "q-graph", remoteQueryGraph)
	osc.AddHandler(server, "q-commands", remoteQueryCommands)
	osc.AddHandler(server, "op", dispatchOperatorCommand)
}



// osc /pig/new optype name
// -> name
// Not used for MIDIInput or MIDIOutput
//
func remoteNewOperator(msg *goosc.Message)([]string, error) {
	template := "ss"
	args, err := osc.Expect(template, osc.ToStringSlice(msg.Arguments))
	if err != nil {
		return empty, err
	}
	otype, name := args[0], args[1]
	op, err := NewOperator(otype, name)
	if err != nil {
		return empty, err
	}
	return osc.StringSlice(op.Name()), err
}

// osc /pig/new-midi-input deviceSpec name
// -> name
//
func remoteNewMIDIInput(msg *goosc.Message)([]string, error) {
	template := "ss"
	args, err := osc.Expect(template, osc.ToStringSlice(msg.Arguments))
	if err != nil {
		return empty, err
	}
	device, name := args[0], args[1]
	op, err := NewMIDIInput(name, device)
	if err != nil {
		return osc.ToStringSlice(msg.Arguments), err
	}
	return osc.StringSlice(op.Name()), err
}

// osc /pig/new-midi-output deviceSpec name
// -> name
//
func remoteNewMIDIOutput(msg *goosc.Message)([]string, error) {
	template := "ss"
	args, err := osc.Expect(template, osc.ToStringSlice(msg.Arguments))
	if err != nil {
		return empty, err
	}
	device, name := args[0], args[1]
	op, err := NewMIDIOutput(name, device)
	if err != nil {
		return osc.ToStringSlice(msg.Arguments), err
	}
	return osc.StringSlice(op.Name()), err
}


// osc /pig/delete-operator name
// -> ACK
//
func remoteDeleteOperator(msg *goosc.Message)([]string, error) {
	var err error
	template := "s"
	args, err := osc.Expect(template, osc.ToStringSlice(msg.Arguments))
	if err != nil {
		return empty, err
	}
	name := args[0]
	_, err = GetOperator(name)
	if err != nil {
		return osc.StringSlice(msg.Arguments), err
	}
	DeleteOperator(name)
	return empty, err
}

// osc /pig/q-operators
// -> list
//
func remoteQueryOperators(msg *goosc.Message)([]string, error) {
	var err error
	ops := Operators()
	acc := make([]string, len(ops))
	for i, op := range ops {
		acc[i] = fmt.Sprintf("%s, %s", op.OperatorType(), op.Name())
	}
	return acc, err
}


// osc /pig/q-roots
// -> list
//
func remoteQueryRoots(msg *goosc.Message)([]string, error) {
	var err error
	ops := RootOperators()
	acc := make([]string, len(ops))
	for i, op := range ops {
		acc[i] = fmt.Sprintf("%s, %s", op.OperatorType(), op.Name())
	}
	return acc, err
}
	
// osc /pig/q-operator-types
// -> list
//
func remoteQueryOperatorTypes(msg *goosc.Message)([]string, error) {
	var err error
	return OperatorTypes(false), err
}

// osc /pig/delete-all-operators
// -> ACK
//
func remoteDeleteAllOperators(msg *goosc.Message)([]string, error) {
	var err error
	ClearRegistry()
	return empty, err
}

// osc /pig/connect parent child
// -> ACK | ERROR
//
func remoteConnect(msg *goosc.Message)([]string, error) {
	var err error
	var parent, child Operator
	template := "ss"
	args, err := osc.Expect(template, osc.ToStringSlice(msg.Arguments))
	if err != nil {
		return empty, err
	}
	parentName, childName := args[0], args[1]
	parent, err = GetOperator(parentName)
	if err != nil {
		return empty, err
	}
	child, err = GetOperator(childName)
	if err != nil {
		return empty, err
	}
	err = parent.Connect(child)
	return empty, err
}
		

// osc /pig/disconnect-child parent child
// -> ACK | ERROR
//
func remoteDisconnect(msg *goosc.Message)([]string, error) {
	var err error
	var parent, child Operator
	template := "ss"
	args, err := osc.Expect(template, osc.ToStringSlice(msg.Arguments))
	if err != nil {
		return empty, err
	}
	parentName, childName := args[0], args[1]
	parent, err = GetOperator(parentName)
	if err != nil {
		return empty, err
	}
	child, err = GetOperator(childName)
	if err != nil {
		return empty, err
	}
	parent.Disconnect(child)
	child.Panic()
	return empty, err
}


// osc /pig/disconnect-all parent
// -> ACK | ERROR
//
func remoteDisconnectAll(msg *goosc.Message)([]string, error) {
	var err error
	var parent Operator
	template := "s"
	args, err := osc.Expect(template, osc.ToStringSlice(msg.Arguments))
	if err != nil {
		return empty, err
	}
	name := args[0]
	parent, err = GetOperator(name)
	if err != nil {
		return empty, err
	}
	parent.DisconnectAll()
	return empty, err
}

// osc /pig/disconnect-parents name
// -> ACK | ERROR
//
func remoteDisconnectParents(msg *goosc.Message)([]string, error) {
	var err error
	var op Operator
	template := "s"
	args, err := osc.Expect(template, osc.ToStringSlice(msg.Arguments))
	if err != nil {
		return empty, err
	}
	name := args[0]
	op, err = GetOperator(name)
	if err != nil {
		return empty, err
	}
	op.DisconnectParents()
	op.Panic()
	return empty, err
}

// osc /pig/q-forest name
// -> list
//
func remoteQueryGraph(msg *goosc.Message)([]string, error) {
	var err error
	var seen = make(map[string]bool)
	var acc []string
	for _, op := range Operators() {
		name := op.Name()
		_, flag := seen[name]
		if !flag {
			seen[name]=true
			if !op.IsLeaf() {
				s := fmt.Sprintf("%s -> [", name)
				for _, c := range op.Children() {
					s += fmt.Sprintf("%s, ", c.Name())
				}
				s = s[0:len(s)-2] + "]"
				acc = append(acc, s)
			} else {
				acc = append(acc, op.Name())
			}
		}
	}
	return acc, err
}
			
		
func remoteQueryCommands(msg *goosc.Message)([]string, error) {
	var err error
	acc := osc.GlobalServer.Commands()
	for _, op := range Operators() {
		name := op.Name()
		for _, cmd := range op.Commands() {
			acc = append(acc, fmt.Sprintf("/pig/op %s, %s", name, cmd))
		}
	}
	return acc, err
}


// osc  /pig/op  [name, command, arguments...]
//
func dispatchOperatorCommand(msg *goosc.Message)([]string, error) {
	var err error
	var args []string
	var op Operator
	template := "ss"
	args, err = osc.Expect(template, osc.ToStringSlice(msg.Arguments))
	if err != nil {
		return empty, err
	}
	name, command := args[0], args[1]
	op, err = GetOperator(name)
	if err != nil {
		return empty, err
	}
	result, rerr := op.DispatchCommand(command, osc.ToStringSlice(msg.Arguments))
	return result, rerr
}
