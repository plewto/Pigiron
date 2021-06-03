package op



import (
	"fmt"
	"strings"
	"errors"
)

// The registry is a global map holding all current operators. 
// MIDIInput and MIDIOutput operators are stored separately.
//
var registry map[string]PigOp = make(map[string]PigOp)

// OperatorExists(name) returns true if the registry contains the named operator.
//
func OperatorExists(name string) bool {
	_, flag := registry[name]
	_, iflag := inputCache[name]
	_, oflag := outputCache[name]
	return flag || iflag || oflag
}

// splitStem separates a string at the final period, if any.
// Ape     --> "Ape" ""
// Ape.Bat --> "Ape" "Bat"
//
func splitStem(s string) (string, string) {
	pos := strings.LastIndex(s, ".")
	head, tail := "", ""
	if pos == -1 {
		head = s
		tail = ""
	} else {
		head = s[:pos]
		tail = s[pos+1:]
	}
	return head, tail
}

// assignName(op) reassigns the operator's name so that it is unique.
// If the registry does not contain an operator by the same name, the
// original name is preserved.   Otherwise the name is modified by
// appending a unique index.
// Returns the actual name.
//
func assignName(op PigOp) string {
	if OperatorExists(op.Name()) {
		base, _ := splitStem(op.Name())
		index := 1
		name := fmt.Sprintf("%v.%d", base, index)
		for OperatorExists(name) {
			index++
			name = fmt.Sprintf("%v.%d", base, index)
		}
		op.setName(name)
	}
	return op.Name()
}

// register adds an operator to the registry.
// If needed, the operator's name is changed to make it unique.
// Returns the actual operator's name.
//
func register(op PigOp) string {
	name := assignName(op)
	registry[name] = op
	return name
}

// DumpRegistry() prints the contents of the operator registry.
//
func DumpRegistry() {
	fmt.Println("Operator registry:")
	for _, op := range registry {
		fmt.Printf("\t%s", op)
	}
}


// MakeOperator(opType, name) creates a new Operator and adds it to the registry.
// All operators should be created by MakeOperator.
//
// opType indicates the type of Operator, valid options are: 
//    Null
//    ...
// name is the proposed name for the operator.  
// The actual name may be different if name is already in use.
//
// Returns the new Operator and an error.
// The error is non nil if opType was invalid.
//
func MakeOperator(opType string, name string) (PigOp, error) {
	var err error
	var op PigOp
	switch opType {
	case "Null":
		op = makeNullOperator(name)
	case "Monitor":
		op  = makeMonitor(name)
	case "ChannelFilter":
		op = makeChannelFilter(name)
	default:
		sfmt := "Invalid Operator type: '%s'"
		msg := fmt.Sprintf(sfmt, opType)
		err = errors.New(msg)
		return op, err
	}
	register(op)
	return op, err
}

// DeleteOperator() removes the named operator from the registry.
// It is not an error if the operator does not exists.
//
func DeleteOperator(name string) {
	delete(registry, name)
}



// ClearRegistry() removes all Operators from the registry.
//
func ClearRegistry() {  // TODO: May want to add op cleanup code?
	for key, _ := range registry {
		delete(registry, key)
	}
}


// GetOperator() returns named operator from the registry.
// An error is returned as the second value and is non nil if no such
// operator exists.
//
func GetOperator(name string) (PigOp, error) {
	var op PigOp
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
func Operators() []PigOp {
	var acc = make([]PigOp, 0, len(registry))
	for _, op := range(registry) {
		acc = append(acc, op)
	}
	for _, op := range inputCache {
		acc = append(acc, op)
	}
	for _, op := range outputCache {
		acc = append(acc, op)
	}
	return acc
}

// RootOperators() returns slice of all root operators.
//
func RootOperators() []PigOp {
	var acc = make([]PigOp, 0, len(registry))
	for _, op := range Operators() {
		if op.IsRoot() {
			acc = append(acc, op)
		}
	}
	return acc
}


