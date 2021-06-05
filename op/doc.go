// Package op defines the primary Pigiron structure, called an Operator.
//
// Operators are defined by the PigOp interface with common behavior
// implemented by the Operator struct.
//
// Each Operator has zero or more parents (inputs) and zero or more
// children (outputs).  Various types of Operators are linked together into
// a "MIDI Process Tree". Cyclical trees are not allowed.
//
// The Operator struct corresponds to an 'abstract class' and is not used
// directly.  Instead several structs extend Operator for specific
// behaviors.
//
// Operators should not be directly constructed.  Use The factory function
// MakeOperator instead.   
//
package op

