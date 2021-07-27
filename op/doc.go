// op package implements the primary Pigiron object, the Operator.
//
// Operators are MIDI processing blocks which may be linked together into
// a MIDI "process-tree".  Each operator has zero or more inputs (called
// it's parents) and zero or more outputs (children).  On reception of a
// MIDI message an operator selectively forwards (a possibly modified)
// version of the message to all of it's children.
//
// The Operator interface defines the basic set of methods.  The
// baseOperator struct provides a concrete implementation of Operator.
//

package op
