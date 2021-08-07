package pattern

import (
	"fmt"
	"testing"
)

func TestCycle(t *testing.T) {
	fmt.Print()
	cy := NewCycle([]int{1, 2, 3})
	v := cy.Value()
	if v != 1 {
		msg := "Cycle test 1, expected intial value of 1, got %d"
		t.Fatalf(msg, v)
	}

	v = cy.Next()
	if v != 1 {
		msg := "Cycle test 2, expected value of 1, got %d"
		t.Fatalf(msg, v)
	}

	v = cy.Next()
	if v != 2 {
		msg := "Cycle test 3, expected value of 2, got %d"
		t.Fatalf(msg, v)
	}

	v = cy.Next()
	if v != 3 {
		msg := "Cycle test 4, expected value of 3, got %d"
		t.Fatalf(msg, v)
	}

	v = cy.Next()
	if v != 1 {
		msg := "Cycle test 5, expected value of 1, got %d"
		t.Fatalf(msg, v)
	}

	cy.Next()
	cy.Reset()
	v = cy.Next()
	if v != 1 {
		msg := "Cycle test 6, expected value of 1 after reset, got %d"
		t.Fatalf(msg, v)
	}
}
