package main

import (
	//"fmt"
	"testing"
	"github.com/plewto/pigiron/op"
)


func TestOperatorBasics(t *testing.T) {
	alpha, _ := op.MakeOperator("Null", "Alpha")
	cmode := alpha.ChannelMode()
	if cmode != op.NoChannel {
		t.Fatalf(".ChannelMode(), expected NoMode, got %v", cmode)
	}
	otype := alpha.OperatorType()
	if otype != "Null" {
		t.Fatalf(".OperatorType(), expected Null, got %s", otype)
	}
	name := alpha.Name()
	if name != "Alpha" {
		t.Fatalf(".Name(), expected Alpha, got %s", name)
	}
	plist := alpha.Parents()
	if len(plist) != 0 {
		t.Fatalf(".Parents(), expected empty list, got %v", plist)
	}
	clist := alpha.Children()
	if len(clist) != 0 {
		t.Fatalf(".Children(), expected empty list, got %v", clist)
	}
	if !alpha.IsRoot() {
		t.Fatalf(".IsRoot() returned false negative")
	}
	if !alpha.IsLeaf() {
		t.Fatalf(".IsLeaf() returned false negative")
	}
}


func TestOperatorConnections(t *testing.T) {
	op.ClearRegistry()
	a, _ := op.MakeOperator("Null", "Alpha")
	b, _ := op.MakeOperator("Null", "Beta")
	c, _ := op.MakeOperator("Null", "Gamma")
	var err error
	
	// simple serial connection
	err = a.Connect(b)
	if err != nil {
		t.Fatalf("Unxepected error from %s.Connect(%s)", "a", "b")
	}
	err = b.Connect(c)
	if err != nil {
		t.Fatalf("Unxepected error from %s.Connect(%s)", "b", "c")
	}

	// circular tree detection
	err = c.Connect(a)
	if err == nil {
		t.Fatalf("Circular tree not detected.")
	}

	// root/leaf test
	var root, leaf bool
	root, leaf = a.IsRoot(), a.IsLeaf()
	if !root || leaf {
		t.Fatalf("a.IsRoot() -> %v   a.IsLeaf() -> %v", root, leaf)
	}

	root, leaf = b.IsRoot(), b.IsLeaf()
	if root || leaf {
		t.Fatalf("b.IsRoot() -> %v   b.IsLeaf() -> %v", root, leaf)
	}

	root, leaf = c.IsRoot(), c.IsLeaf()
	if root || !leaf {
		t.Fatalf("c.IsRoot() -> %v   c.IsLeaf() -> %v", root, leaf)
	}

	if !(b.IsChildOf(a) && c.IsChildOf(b)) {
		t.Fatalf("IsChildOf returned false negative")
	}

	if !(a.IsParentOf(b) && b.IsParentOf(c)) {
		t.Fatalf("IsParentOf returned false negative")
	}
	
	a.Disconnect(b)
	root, leaf = a.IsRoot(), a.IsLeaf()
	if !root || !leaf {
		t.Fatalf("a.isRoot() -> %v,  a.isLeaf() -> %v, after a.Disconnect(b)", root, leaf)
	}

	root, leaf = b.IsRoot(), b.IsLeaf()
	if !root || leaf {
		t.Fatalf("b.isRoot() -> %v,  b.isLeaf() -> %v, after a.Disconnect(b)", root, leaf)
	}


	

	roots := op.RootOperators()
	if len(roots) != 2 {
		t.Fatalf("Expected len(roots) of 2, got %d", len(roots))
	}

}
	
			

	

