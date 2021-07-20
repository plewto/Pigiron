package midi

import (
	"testing"
	"fmt"
	"github.com/plewto/pigiron/pigpath"
)


func TestReadSMF(t *testing.T) {
	fmt.Print()
	filename := pigpath.ResourceFilename("testFiles", "a.mid")
	fmt.Printf("Test filename: %s\n", filename)
	smf, err := ReadSMF(filename)
	if err != nil {
		errmsg := "\nReadSMF returned error for known good MIDI file: %s"
		errmsg += "\n%s"
		t.Fatalf(errmsg, filename, err)
	}
	smf.Dump()
}
