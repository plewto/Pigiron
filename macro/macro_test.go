package macro

import (
	"fmt"
	"testing"
)


func TestMacro(t *testing.T) {
	fmt.Print()
	Define("Alpha", "Apple", []string{"A", "$0", "$1"})
	if !IsMacro("Alpha") {
		msg := "Macro test 1, Expected IsMacro() to return true, got false"
		t.Fatalf(msg)
	}
	if IsMacro("NOT-DEFINED") {
		msg := "Macro test 2, Expected IsMacro() to return false, got true"
		t.Fatalf(msg)
	}

	expect := "Apple A, Bat, Cat"
	ex, err := Expand("Alpha", []string{"Bat", "Cat"})
	if err != nil {
		msg := "Macro test 3, got unexpected error: %s"
		t.Fatalf(msg, err)
	}
	if ex != expect {
		msg := "Macro test 4, Expand returned '%s', expected '%s'"
		t.Fatalf(msg, ex, expect)
	}

	Delete("Alpha")
	if IsMacro("Alpha") {
		msg := "Macro test 5, IsMacro returns true after macro was deleted"
		t.Fatalf(msg)
	}
}


