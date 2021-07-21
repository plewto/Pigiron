package midi

import (
	"testing"
	"fmt"
	"github.com/plewto/pigiron/pigpath"
)


func TestReadSMF(t *testing.T) {

	fmt.Println("*** EXPECT TO SEE WARNINGS ***")


	formatFilename := func(name string) string {
		filename := pigpath.ResourceFilename("testFiles", name)
		fmt.Printf("Using test file '%s'\n", filename)
		return filename
	}

	filename := formatFilename("a1.mid")
	smf, err := ReadSMF(filename)
	if err != nil {
		errmsg := "\nReadSMF returnd error for known good MIDI file %s\n"
		errmsg += "%s\n"
		t.Fatalf(errmsg, filename, err)
	}

	if smf.TrackCount() != 1 {
		errmsg := "Expected track count = 1, got %d"
		t.Fatalf(errmsg, smf.TrackCount())
	}
	_, err = smf.Track(0)
	if err != nil {
		errmsg := "\nsmf.Track(0) returned unexpected error"
		errmsg += "\n%s"
		t.Fatalf(errmsg, err)
	}
	_, err = smf.Track(19)
	if err == nil {
		errmsg := "\nsmf.Track(19) did not detect track-number error"
		t.Fatalf(errmsg)
	}

			
	filename = formatFilename("b1.mid")
	_, err = ReadSMF(filename)
	if err == nil {
		errmsg := "\nreadSMF(%s) did not return error for non-midi file"
		t.Fatalf(errmsg)
	}



}
