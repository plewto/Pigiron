package op

import (
	"fmt"
	
	"github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/smf"
)


// Implements Operator and Transport interfaces

type MIDIPlayer struct {
	baseOperator
	smf *smf.SMF
	noteQueue *midi.NoteQueue
	isPlaying bool
	currentTime float64
	
}

func newMIDIPlayer(name string) *MIDIPlayer {
	op := new(MIDIPlayer)
	initOperator(&op.baseOperator, "MIDIPlayer", name, midi.NoChannel)
	op.noteQueue = midi.MakeNoteQueue()
	op.isPlaying = false
	op.currentTime = 0.0
	return op
}

func (op *MIDIPlayer) Info() string {
	acc := op.commonInfo()
	acc += fmt.Sprintf("\tfilename  : \"%s\"\n", op.MediaFilename())
	acc += fmt.Sprintf("\tIsPlaying : %v\n", op.IsPlaying())
	return acc
}
	
func (op *MIDIPlayer) Reset() {
	op.Stop()
}

func (op *MIDIPlayer) Panic() { // TODO: transmit all notes off
	op.Stop()
	
}

func (op *MIDIPlayer) Stop() { // TODO Implement
}

func (op *MIDIPlayer) Play() { // TODO Implement
}

func (op *MIDIPlayer) Continue() { // TODO Implement
}

func (op *MIDIPlayer) IsPlaying() bool { // TODO Implement
	return false
}

func (op *MIDIPlayer) LoasMedia(filename string) error { // TODO Implement
	var err error
	return err
}

func (op *MIDIPlayer) MediaFilename() string { 
	if op.smf != nil {
		return op.smf.Filename()
	} else {
		return ""
	}
}

