package pattern

import (
	"fmt"
)


/*
** Pattern interface defines generalized sequendce generator.
**
** Value() int returns current Pattern value.
** SetValues([]int) sets pattern values.
** Next() returns the next pattern value.  
** 
*/ 
type Pattern interface {
	Value() int
	SetValues([]int)
	Next() int
	Reset()
	String()
}


/*
** Cycle struct implements the Pattern interface
** Cycle values are returned in a cyclical manner.
**
*/
type Cycle struct {
	values []int
	pointer int
}

// NewCycle(v1, v2, v3 ...int)
// Returns new pointer to Cycle with given values.
//
func NewCycle(values []int) *Cycle {
	c := new(Cycle)
	c.SetValues(values)
	return c
}
	
func (cy *Cycle) Value() int {
	return cy.values[cy.pointer]
}

func (cy *Cycle) SetValues(values []int) {
	cy.values = values
	cy.Reset()
}
	

func (cy *Cycle) Next() int {
	v := cy.Value()
	cy.pointer = (cy.pointer + 1) % len(cy.values)
	return v
}

func (cy *Cycle) Reset() {
	cy.pointer = 0
}

func (cy *Cycle) String() string {
	acc := "Cycle: "
	for _, v := range cy.values {
		acc += fmt.Sprintf("%d ", v)
	}
	return acc
}
	
