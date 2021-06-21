package osc

import (
	"fmt"
	"sort"
	"github.com/plewto/pigiron/op"
)


var oscHelp = make(map[string]func())

func init() {
	oscHelp["help"] = helpHelp
	oscHelp["ping"] = pingHelp
	oscHelp["exit"] = exitHelp
	oscHelp["new-operator"] = newOperatorHelp
	oscHelp["new-midi-input"] = newMIDIInputHelp
	oscHelp["new-midi-output"] = newMIDIOutputHelp
	oscHelp["delete-operator"] = deleteOperatorHelp
	oscHelp["connect"] = connectHelp
	oscHelp["disconnect"] = disconnectHelp
	oscHelp["disconnect-all"] = disconnectAllHelp
	oscHelp["destroy-forest"] = destroyForestHelp
	oscHelp["print-forest"] = printForestHelp
	oscHelp["q-is-parent"] = queryIsParentHelp
	oscHelp["q-midi-inputs"] = queryMIDIInputsHelp
	oscHelp["q-midi-outputs"] = queryMIDIInputsHelp
	oscHelp["q-operators"] = queryOperatorsHelp
	oscHelp["q-roots"] = queryRootsHelp
	oscHelp["q-children"] = queryChildrenHelp
	oscHelp["q-parents"] = queryParentsHelp
	oscHelp["panic"] = panicHelp
	oscHelp["reset"] = resetHelp
	oscHelp["batch"] = batchHelp
}

func helpHelp() {
	fmt.Println("\nHelp topics:")
	keys := make([]string, 0, len(oscHelp))
	for key, _ := range oscHelp {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Printf("\t%s\n", key)
	}
	prompt()
	
}

func helpPrefix(topic string, arguments string) {
	fmt.Printf("\n%s help\n", topic)
	fmt.Printf("\tOSC command: /pig/%s %s\n", topic, arguments)
}

func helpReturn(returns string) {
	fmt.Printf("\treturns: %s\n", returns)
	prompt()
}

func pingHelp() {
	helpPrefix("ping", "")
	fmt.Println("\tDiagnostic, prints \"ping\".")
	helpReturn("ACK")
}

func exitHelp() {
	helpPrefix("exit", "")
	fmt.Println("\tTerminates application.")
	helpReturn("Never returns.")
}

func newOperatorHelp() {
	helpPrefix("new-operator", "operator-type, operator-name")
	fmt.Println("\tCreates new operator")
	fmt.Println("\toperator-type may be one of:")
	for _, opt := range op.OperatorTypes(true) {
		fmt.Printf("\t\t%s\n", opt)
	}
	fmt.Println("\tYou can not use new-operator for MIDIInput or MIDIOutput.")
	fmt.Println("\tThe actual operator's name may be modified.")
	helpReturn("ACK actual-name")
}
	
func newMIDIInputHelp() {
	helpPrefix("new-midi-input", "device, operator-name")
	fmt.Println("\tCreates new MIDIInput operator.")
	fmt.Println("\tUse command 'q-midi-inputs' for list of available devices.")
	fmt.Println("\tThe device argument may be a sub-string of the full device name.")
	fmt.Println("\tThe actual operator's name may be modified.")
	helpReturn("ACK actual-name")
}

func newMIDIOutputHelp() {
	helpPrefix("new-midi-output", "device, operator-name")
	fmt.Println("\tCreates new MIDIOutput operator.")
	fmt.Println("\tUse command 'q-midi-outputs' for list of available devices.")
	fmt.Println("\tThe device argument may be a sub-string of the full device name.")
	fmt.Println("\tThe actual operator's name may be modified.")
	helpReturn("ACK actual-name")
}

func deleteOperatorHelp() {
	helpPrefix("delete-operator", "operator-name")
	fmt.Println("\tDeletes named operator.")
	fmt.Println("\tIt is not an error if the operator does not exists.")
	helpReturn("ACK")
}
	
func connectHelp() {
	helpPrefix("connect", "parent-name, child-name")
	fmt.Println("\tConnects parent as MIDI input of child.")
	fmt.Println("\tIt is an error if the connection causes a circular tree.")
	helpReturn("ACK")
}


func disconnectHelp() {
	helpPrefix("disconnect", "parent-name, child-name")
	fmt.Println("\tRemoves parent as MIDI input to child.")
	fmt.Println("\tIt is not an error if the operators are not currently connected.")
	helpReturn("ACK")
}

func disconnectAllHelp() {
	helpPrefix("disconnect-all", "parent-name")
	fmt.Println("\tDisconnect all child operators from parent.")
	helpReturn("ACK")
}
		
func destroyForestHelp() {
	helpPrefix("destroy-forest", "")
	fmt.Println("\tRemove connections from ALL operators.")
	helpReturn("ACK")
}

func printForestHelp() {
	helpPrefix("print-forest", "")
	fmt.Println("\tPrints representation for all operator connections.")
	helpReturn("ACK")
}

func queryIsParentHelp() {
	helpPrefix("q-is-parent", "parent-name, child-name")
	helpReturn("ACK true if parent is a parent to child.")
}

func queryMIDIInputsHelp() {
	helpPrefix("q-midi-inputs", "")
	helpReturn("ACK list of MIDI input devices.")
}

func queryMIDIOutputsHelp() {
	helpPrefix("q-midi-outputs", "")
	helpReturn("ACK list of MIDI output devices.")
}

func queryOperatorsHelp() {
	helpPrefix("q-operators", "")
	helpReturn("ACK list of all operator names.")
}

func queryRootsHelp() {
	helpPrefix("q-roots", "")
	helpReturn("ACK list of all root operator names.")
}

func queryChildrenHelp() {
	helpPrefix("q-children", "operator-name")
	helpReturn("ACK list of all operator's children")
}

func queryParentsHelp() {
	helpPrefix("q-parents", "operator-name")
	helpReturn("ACK list of all of operator's parents")
}	

func panicHelp() {
	helpPrefix("panic", "")
	fmt.Println("\tHalts playback and kills all notes.")
	helpReturn("ACK")
}

func resetHelp() {
	helpPrefix("reset", "")
	fmt.Println("\tResets all operators")
	helpPrefix("reset/op/<operator-name>", "")
	fmt.Println("\tResets named operator")
	helpReturn("ACK")
}


func batchHelp() {
	fmt.Printf("\nbatch help\n")
	fmt.Printf("\tSpecial command:  batch filename\n")
	fmt.Printf("\tReads commands from file.\n")
	fmt.Printf("\tThe contents of the file is read line-by-line and treated\n")
	fmt.Printf("\tas if each line was enterd interactivly.\n")
	fmt.Printf("\tLine which begin with # are ignored, there must be a space after #.\n")
	fmt.Printf("\tfilename may being with '~/' to indicate the User's home directory.\n")
	prompt()
}

		
