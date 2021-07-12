package smf

import (
	"fmt"
	"os"
)

func exError(message string, more ...string) error {
	msg := fmt.Sprintf("ERROR: %s", message)
	for _, s := range more {
		msg += fmt.Sprintf("\nERROR: %s", s)
	}
	return fmt.Errorf(msg)
}

func compoundError(err1 error, message string, more ...string) error {
	msg := fmt.Sprintf("ERROR: %s\n", message)
	for _, s := range more {
		msg += fmt.Sprintf("\nERROR: %s", s)
	}
	msg += fmt.Sprintf("PREVIOUS ERROR:\n%s", err1)
	return fmt.Errorf(msg)
}


// SMF struct Represents a Standard MIDI File.
//
type SMF struct {
	filename string
	header *Header
	tracks []*Track
}


func ReadSMF(filename string) (*SMF, error) {
	var err error
	var smf = &SMF{}
	smf.filename = filename
	var header *Header
	var file *os.File
	file, err = os.Open(filename)
	errmsg := fmt.Sprintf("smf.ReadSMF, can not open file '%s'", filename)
	if err != nil {
		err = compoundError(err, errmsg)
		return smf, err
	}
	defer file.Close()
	header, err = readHeader(file)
	if err != nil {
		err = compoundError(err, errmsg)
		return smf, err
	}
	smf.header = header
	for i := 0; i < header.chunkCount; i++ {
		var id [4]byte
		var length int
		id, length, err = readChunkPreamble(file)
		if err != nil {
			err = compoundError(err, errmsg)
			return smf, err
		}
		if id != trackID { // skip non-track chunks
			msg := "Skiping non-recognized chunk %v, smf file '%s'\n"
			fmt.Printf(msg, string(id[:]), filename)
			var junk = make([]byte, length)
			_, err = file.Read(junk)
			if err != nil {
				sid := string(id[:])
				msg = fmt.Sprintf("smf.ReadSMF file: '%s'", filename)
				msg2 := "An error occured while reading non-recognized chunk: %s"
				err = compoundError(err, msg, fmt.Sprintf(msg2, sid))
				return smf, err
			}
		} else { // read track chunk
			trk := &Track{make([]byte, 0, 1024)}
			trk.bytes, err = readTrackBytes(file, length)
			if err != nil {
				msg := fmt.Sprintf("smf.READSMF file: '%s'", filename)
				msg2 := "An error occured while reading track chunk %d"
				err = compoundError(err, msg, fmt.Sprintf(msg2, i))
				return smf, err
			}
			smf.tracks = append(smf.tracks, trk)
		}
	}
	return smf, err
}
				
				
	
	
