package smf

import (
	"fmt"
	"testing"
)

var (
	GOOD_FILES = []string{"!/resources/testFiles/a1.mid",
		"!/resources/testFiles/a2.mid"}
	BAD_FILES = []string{"!/resources/testFiles/b1.mid",
		"!/resources/testFiles/b2.mid",
		"!/resources/testFiles/b3.mid"}
	RECOVERABLE_FILES = []string{"!/resources/testFiles/c1.mid",
		"!/resources/testFiles/c2.mid"}
)


func TestSMF(t *testing.T) {
	fmt.Println("TestSMF")
	for _, f := range GOOD_FILES {
		smf, err := ReadSMF(f)
		if err != nil {
			errmsg := "Unexpected error while reading smf file '%s'\n%s"
			t.Fatalf(errmsg, f, err.Error())
		}
		fmt.Printf("test file: %s\n", smf.filename)
	}
	for _, f := range BAD_FILES {
		_, err := ReadSMF(f)
		if err == nil {
			errmsg := "Did not detect non-smf file: %s"
			t.Fatalf(errmsg, f)
		}
	}
	fmt.Println("*** EXPECT TO SEE WARNINGS ***")
	for _, f := range RECOVERABLE_FILES {
		_, err := ReadSMF(f)
		if err != nil {
			errmsg := "Did not read recoverable SMF file: %s\n%s"
			t.Fatalf(errmsg, f, err.Error())
		}
	}
}
