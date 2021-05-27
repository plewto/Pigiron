package op

// NullOperator is an Operator with no additional behavior.
// Its sole purpose is for testing.
// Do not construct NullOperator directly, instead use the registry
// MakeOperator factory function.
//
type NullOperator struct {
	Operator
}

func makeNullOperator(name string) *NullOperator {
	op := new(NullOperator)
	initOperator(&op.Operator, "Null", name, NoChannel)
	return op
}
	
	
