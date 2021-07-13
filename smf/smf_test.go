package smf

import (
	"testing"
	"fmt"

	"github.com/plewto/pigiron/fileio"
	
)

func TestReadSMF(t *testing.T) {
	if noResourcesAbort("TestReadSMF") {
		return
	}
	var err error
	var smf *SMF
	filename, _ := fileio.ResourceFilename("testFiles", "a.mid")
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
	if smf.header.division != 24 {
		msg := "smf header does not have expected division.\n"
		msg += "Expected 24, got %d"
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


// Test various malformed header chunks.
func TestReadSMFJunk (t *testing.T) {
	if noResourcesAbort("TestReadSMFJunk") {
		return
	}
	var err error

	// File does not exists
	filename, _ := fileio.ResourceFilename("testFiles", "does-not-exists.mid")
	_, err = ReadSMF(filename)
	if err == nil {
		msg := "\nReadSMF did not return an error for non-existent file\n"
		msg += fmt.Sprintf("filename was %s", filename)
		t.Fatal(msg)
	}

	// File exists but has malformed header id.
	filename, _ = fileio.ResourceFilename("testFiles", "bad1.mid")
	_, err = ReadSMF(filename)
	if err == nil {
		msg := "\nReadSMF did not return an error for badly malformed MIDI file\n"
		msg += fmt.Sprintf("filename was %s", filename)
		t.Fatalf(msg)
	}

	// unsupported format test
	filename, _ = fileio.ResourceFilename("testFiles", "badFormat.mid")
	_, err = ReadSMF(filename)
	if err == nil {
		msg := "\nReadSMF did not return an error for unsupported format\n"
		msg += fmt.Sprintf("filename was %s", filename)
		t.Fatal(msg)
	}

	// no tracks
	filename, _ = fileio.ResourceFilename("testFiles", "badNoTracks.mid")
	_, err = ReadSMF(filename)
	if err == nil {
		msg := "\nReadSMF did not return an error for zero track count\n"
		msg += fmt.Sprintf("filename was %s", filename)
		t.Fatal(msg)
	}

	// weird division
	filename, _ = fileio.ResourceFilename("testFiles", "badWeirdDivision.mid")
	_, err = ReadSMF(filename)
	if err == nil {
		msg := "\nReadSMF did not return an error weird looking clock division"
		msg += fmt.Sprintf("filename was %s", filename)
		t.Fatal(msg)
	}
	
}
	
