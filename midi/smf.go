package midi

import (
	"fmt"
	"os"
	"github.com/plewto/pigiron/pigerr"
)

type SMF struct {
	filename string
	header *SMFHeader
	tracks[] *SMFTrack
}

func (smf *SMF) String() string {
	return fmt.Sprintf("SMF '%s'", smf.filename)
}

func (smf *SMF) Dump() {
	fmt.Println("SMF")
	fmt.Printf("  filename : %s\n", smf.filename)
	fmt.Printf("  format   : %d\n", smf.header.format)
	fmt.Printf("  tracks   : %d\n", smf.header.trackCount)
	fmt.Printf("  division : %d\n", smf.header.division)
	for i, trk := range smf.tracks {
		fmt.Printf("  TRACK: %d  events: %d\n", i, trk.Length())
		for j, evnt := range trk.events {
			fmt.Printf("  [%3d] %s\n", j, evnt)
		}
	}
}

func ReadSMF(filename string) (smf *SMF, err error) {
	var file *os.File
	file, err = os.Open(filename)
	if err != nil {
		errmsg := "Can not open MIDI file, filename = %s"
		err = pigerr.CompoundError(err, fmt.Sprintf(errmsg, filename))
		return
	}
	defer file.Close()
	smf = &SMF{}
	smf.filename = filename
	smf.header, err = readSMFHeader(file)
	// Due to the possible presence of non-track chunks, the
	// header.trackCount may initially be too high.
	// The count is adjusted once all chunks have been read.
	// non-track chunks are not-supported and are discarded.
	if err != nil {
		errmsg := "MIDI file header mallformed, filename =  %s"
		err = pigerr.CompoundError(err, fmt.Sprintf(errmsg, filename))
		return
	}
	smf.tracks = make([]*SMFTrack, 0, smf.header.trackCount)
	for i := 0; i < smf.header.trackCount; i++ {
		var id chunkID
		var bytes []byte
		var track *SMFTrack
		id, bytes, err = readRawChunk(file)
		if err != nil {
			errmsg := "Error while reading MIDI file chunk %d, filename = %s"
			err = pigerr.CompoundError(err, fmt.Sprintf(errmsg, i, filename))
			return
		}
		if !id.eq(trackID) {
			errmsg := "Ignoring non-track chunk number %d, type %s, filename = %s"
			pigerr.Warning(fmt.Sprintf(errmsg, i, id, filename))
			continue
		}
		track, err = convertTrackBytes(bytes)
		if err != nil {
			errmsg := "Error wile converting track %d, filename = %s"
			err = pigerr.CompoundError(err, fmt.Sprintf(errmsg, i, filename))
			return
		}
		smf.tracks = append(smf.tracks, track)
	}
	smf.header.trackCount = len(smf.tracks)
	return
}

		
	
	

