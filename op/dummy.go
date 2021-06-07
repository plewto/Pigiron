package op

import midi "github.com/plewto/pigiron/midi"

// DummyOperator is an Operator with no additional behavior.
// Its sole purpose is for testing.
//
// Do not construct DummyOperator directly, instead use the
// MakeOperator factory function.
//
type DummyOperator struct {
	Operator
}

func makeDummyOperator(name string) *DummyOperator {
	op := new(DummyOperator)
	initOperator(&op.Operator, "Dummy", name, midi.NoChannel)
	return op
}
	
	
