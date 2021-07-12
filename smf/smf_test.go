package smf

import (
	"testing"
	"fmt"

	"github.com/plewto/pigiron/fileio"
	
)


func TestReadSMF(t *testing.T) {
	filename, err := fileio.ResourceFilename("testFiles", "a.mid")
	if err != nil {
		fmt.Println("\nWARNING: Can not read resource file required for TestReadSMF")
		fmt.Println("WARNING: Error from fileio.ResourceFilename was:")
		fmt.Printf("%s\n", err)
		fmt.Printf("WARNING: Aborting test.\n\n")
		return
	}
	var smf *SMF
	smf, err = ReadSMF(filename)
	if err != nil {
		msg := "smf.ReadSMF returned unexpected error for file %s\n"
		err = compoundError(err, fmt.Sprintf(msg, filename))
		t.Fatal(err)
	}
	if smf.filename != filename {
		msg := "*SMF.filename does not equal test filename\n '%s' != '%s'\n"
		t.Fatalf(msg, smf.filename, filename)
	}
	if smf.header.ID() != headerID {
		msg := "smf header does not have expected id.\n"
		msg += "Expected %s, got %s\n"
		t.Fatalf(msg, headerID, smf.header.ID())
	}
	if smf.header.chunkCount != 1 {
		msg := "smf header does not have expected chunkCount.\n"
		msg += "Expected 1, got %d"
		t.Fatalf(msg, smf.header.chunkCount)
	}
	if smf.header.division != 480 {
		msg := "smf header does not have expected division.\n"
		msg += "Expected 480, got %d"
		t.Fatalf(msg, smf.header.division)
	}
	// check consistency between header chunkCount and actual track count.
	// Note non-recognized chucks are discarded.
	if smf.header.chunkCount != len(smf.tracks) {
		msg := "Inconsistency between header chunkCount and actual track count.\n"
		msg += "The header chunkCount should ignore non-track chunks.\n"
		msg += "header chunkCount = %d, actual track count = %d"
		t.Fatalf(msg, smf.header.chunkCount, len(smf.tracks))
	}
	// Check track contents
	trk := smf.tracks[0]
	if trk.ID() != trackID {
		msg := "smf track does not have expected id.\n"
		msg += "Expected %s, got %s\n"
		t.Fatalf(msg, trackID, trk.ID())
	}
	if trk.Length() != len(trk.Bytes()) {
		msg := "Inconsistency between track length value and actual byte count.\n"
		msg += "track.Length() = %d, len(trk.Bytes()) = %d"
		t.Fatalf(msg, trk.Length(), len(trk.Bytes()))
	}

	
}
