package op

import (
	"fmt"
	
	"github.com/plewto/pigiron/midi"
)


// Implements Operator and Transport interfaces

type MIDIPlayer struct {
	baseOperator
	currentFilename string
}

func newMIDIPlayer(name string) *MIDIPlayer {
	op := new(MIDIPlayer)
	initOperator(&op.baseOperator, "MIDIPlayer", name, midi.NoChannel)
	op.currentFilename = ""
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
	op.currentFilename = ""
}

func (op *MIDIPlayer) Panic() {
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

func (op *MIDIPlayer) MediaFilename() string { // TODO Implement
	return op.currentFilename
}

