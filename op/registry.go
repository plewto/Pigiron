package op

import (
	"time"
	"fmt"
	"errors"
)

var OperatorTypes = []string{
	"ChannelFilter",
	"Disrtributor",
	"MIDIInput",
	"MIDIOutput",
	"MIDIPlayer",
	"Monitor",
	"Transformer"}

// The registry is a global map holding all current operators. 
// MIDIInput and MIDIOutput operators are stored separately.
//
var registry map[string]Operator = make(map[string]Operator)


// OperatorExists(name) returns true if the registry contains the named operator.
//
func OperatorExists(name string) bool {
	_, flag := registry[name]
	return flag
}


// register() adds an operator to the registry.
// Returns operator's name
//
func register(op Operator) string {
	registry[op.Name()] = op
	return op.Name()
}


// DumpRegistry() prints the contents of the operator registry.
//
func DumpRegistry() {
	fmt.Println("Operator registry:")
	for _, op := range registry {
		fmt.Printf("\t%s", op)
	}
}


// NewOperator function creates a new Operator and adds it to the registry.
// All operators should be created by NewOperator.
//
// opType indicates the type of Operator
//
// Returns the new Operator and an error.
// The error is non-nil if:
//     1) opType was invalid
//     2) if an Operator name exists and its type does not match opType.
//
// If Operator name exist with the same type as opType, the existing
// operator is reused.
//
func NewOperator(opType string, name string) (Operator, error) {
	var err error
	var op Operator
	if OperatorExists(name) {
		other, _ := registry[name]
		if other.OperatorType() != opType {
			msg := "An operator named %s of type %s already exists\n"
			msg += "Can not create new %s Operator with same name."
			err = fmt.Errorf(msg, name, other.OperatorType(), opType)
			return op, err
		} else {
			return GetOperator(name)
		}
	}
	switch opType {
	case "Dummy":
		op = newDummyOperator(name)
	case "Monitor":
		op  = newMonitor(name)
	case "ChannelFilter":
		op = newChannelFilter(name)
	case "Distributor":
		op = newDistributor(name)
	case "MIDIPlayer":
		op = newMIDIPlayer(name)
	case "Transformer":
		op = newTransformer(name)
	default:
		sfmt := "Invalid Operator type: '%s'"
		msg := fmt.Sprintf(sfmt, opType)
		err = errors.New(msg)
		return op, err
	}
	register(op)
	return op, err
}

// DeleteOperator() Deletes named operator.
// Returns error if operator does not exists or it is a MIDIInput.
//
func DeleteOperator(name string) error {
	var err error
	var op Operator
	op, err = GetOperator(name) 
	if err != nil {
		return err
	}
	if op.OperatorType() == "MIDIInput" {
		msg := "Can not delete MIDIInput Operator: %s"
		err = fmt.Errorf(msg, name)
		return err
	}
	op.Panic()
	time.Sleep(1 * time.Millisecond)
	op.DisconnectAll()
	op.Close()
	delete(registry, name)
	return err
}



// ClearRegistry() Deletes all Operators (except MIDIInputs)
//
func ClearRegistry() {
	for _, root := range RootOperators() {
		root.Panic()
	}
	time.Sleep(10 * time.Millisecond)
	for _, op := range Operators() {
		op.DisconnectAll()
		if op.OperatorType() != "MIDIInput" {
			delete(registry, op.Name())
			op.Close()
		}
	}
}


// GetOperator() returns named operator.
// Returns:
//   1. The operator
//   2. non-nil error if the operator does not exists.
//
func GetOperator(name string) (Operator, error) {
	var op Operator
	var err error
	if OperatorExists(name) {
		op = registry[name]
	} else {
		sfmt := "Operator '%s' does not exists"
		msg := fmt.Sprintf(sfmt, name)
		err = errors.New(msg)
	}
	return op, err
}


// Operators() returns unordered slice of all current operators.
// 
func Operators() []Operator {
	var acc = make([]Operator, 0, len(registry))
	for _, op := range(registry) {
		acc = append(acc, op)
	}
	return acc
}

// RootOperators() returns slice of all root operators.
//
func RootOperators() []Operator {
	var acc = make([]Operator, 0, len(registry))
	for _, op := range Operators() {
		if op.IsRoot() {
			acc = append(acc, op)
		}
	}
	return acc
}


func DestroyForest() {
	for _, root := range RootOperators() {
		root.DisconnectTree()
	}
}

func Cleanup() {
	for _, op := range registry {
		op.Close()
	}
}


func ResetAll() {
	for _, op := range Operators() {
		op.Reset()
	}
}

func PanicAll() {
	for _, op := range RootOperators() {
		op.Panic()
	}
}
