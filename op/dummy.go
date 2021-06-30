package op

import midi "github.com/plewto/pigiron/midi"

// DummyOperator is an Operator with no additional behavior.
// Its sole purpose is for testing.
//
type DummyOperator struct {
	baseOperator
}

func newDummyOperator(name string) *DummyOperator {
	op := new(DummyOperator)
	initOperator(&op.baseOperator, "Dummy", name, midi.NoChannel)
	return op
}
	
	
