package main

import (
	//"fmt"
	"testing"
	"github.com/plewto/pigiron/op"
)


func TestRegistry(t *testing.T) {
	op.ClearRegistry()

	// Capture incorrect Operator type.
	_, err := op.MakeOperator("Wrong", "Wrong")
	if err == nil {
		t.Fatal("MakeOperator() did not flag an error for invalid Operator type.")
	}

	// Check for false error with valid Operator type.
	name, name1, name2 := "Zeta", "Zeta.1", "Zeta.2"
	_, err = op.MakeOperator("Null", name)
	if err != nil {
		t.Fatalf("MakeOperator(%s, %s) returnd incorrect error.", "Null", name)
	}

	if !op.OperatorExists(name) {
		t.Fatalf("OperatorExists(%s) returned false negative.", name)
	}

	var zop op.PigOp
	zop, err = op.GetOperator(name)
	if err != nil {
		t.Fatalf("GetOperator(%s) retruned incorrect error.", name)
	}

	// Chack unique Operator name was not mangled.
	if zop.Name() != name {
		t.Fatalf("zop = GetOperator(%s), expected zop.Name() %s, got %s", name, name, zop.Name())
	}

	// Check auto name-mangling
	zop1, _ := op.MakeOperator("Null", name)
	zop2, _ := op.MakeOperator("Null", name)
	if zop1.Name() != name1 || zop2.Name() != name2 {
		t.Fatalf("Name mangling expected %s and %s, got %s, %s", name1, name2, zop1.Name(), zop2.Name())
	}
	
	
	oplist := op.Operators()
	if len(oplist) != 3 {
		t.Fatalf("Operators(), expected list of length 1, found %v", oplist)
	}

	op.DeleteOperator(name1)
	if op.OperatorExists(name1) {
		t.Fatalf("OperatorExists(%s) returned false positive after DeleteOperator(%s)", name1, name1)
	}

	op.ClearRegistry()
	if len(op.Operators()) != 0 {
		t.Fatalf("registry not empty after ClearRegistry()")
	}
}


