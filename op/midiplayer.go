package op

import (
	"fmt"
	"time"
	"github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/pigerr"
	"github.com/plewto/pigiron/pigpath"
)


// Implements Operator and Transport interfaces

type MIDIPlayer struct {
	baseOperator
	smf *midi.SMF
	noteQueue *midi.NoteQueue
	isPlaying bool
	eventIndex int
	tempo float64
	tickDuration float64
	currentTime int // msec
	delayStart int  // msec
	
}

func newMIDIPlayer(name string) *MIDIPlayer {
	op := new(MIDIPlayer)
	initOperator(&op.baseOperator, "MIDIPlayer", name, midi.NoChannel)
	op.noteQueue = midi.MakeNoteQueue()
	op.isPlaying = false
	op.eventIndex = 0
	op.tempo = 120.0
	op.tickDuration = midi.TickDuration(24, op.tempo)
	op.currentTime = 0.0  // milliseconds
	op.delayStart = 200   // milliseconds
	initTransportHandlers(op)
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
	op.eventIndex = 0
	op.tempo = 120.0
	if op.smf != nil {
		op.tickDuration = midi.TickDuration(op.smf.Division(), op.tempo)
	} else {
		op.tickDuration = midi.TickDuration(24, op.tempo)
	}
	op.currentTime = 0.0
}

func (op *MIDIPlayer) Panic() {
	op.Reset()
}

func (op *MIDIPlayer) Stop() {
	fmt.Printf("\nMIDIPlayer %s Stop()\n", op.Name())
	op.isPlaying = false
	for ci := 0; ci < 16; ci++ {
		ch := midi.MIDIChannelNibble(ci)
		ev := midi.MakeControllerEvent(ch, 120, 127)
		op.distribute(ev.PortmidiEvent())
		time.Sleep(1 * time.Millisecond)
		ev = midi.MakeControllerEvent(ch, 120, 0)
		op.distribute(ev.PortmidiEvent())
	}
	op.killActiveNotes()
}

func (op *MIDIPlayer) killActiveNotes() {
	fmt.Printf("MIDIPlayer %s killActiveNotes\n", op.Name())
	counter := 0
	for _, pme := range op.noteQueue.OffEvents() {
		op.distribute(*pme)
		if counter % 16 == 0 {
			time.Sleep(2 * time.Millisecond)
		}
		counter++
	}
	op.noteQueue.Reset()
}
	

func (op *MIDIPlayer) Play() error {
	fmt.Printf("MIDIPlayer %s Play()\n", op.Name())
	var err error
	if op.isPlaying {
		return err
	}
	op.Stop()
	if op.smf == nil {
		errmsg := "No MIDI file loaded"
		err = pigerr.New(errmsg)
		return err
	}
	op.eventIndex = 0
	op.tempo = 120.0
	op.tickDuration = midi.TickDuration(op.smf.Division(), op.tempo)
	op.currentTime = 0.0
	op.isPlaying = true
	go op.playloop()
	return err
}

func (op *MIDIPlayer) Continue() error {
	fmt.Printf("MIDIPlayer %s Continue()", op.Name())
	var err error
	if op.isPlaying {
		return err
	}
	if op.smf == nil {
		errmsg := "NO MIDI file loaded"
		return pigerr.New(errmsg)
	}
	op.isPlaying = true
	go op.playloop()
	return err
}

func (op *MIDIPlayer) IsPlaying() bool {
	return op.isPlaying
}

func (op *MIDIPlayer) LoadMedia(filename string) error {
	var err error
	var smf *midi.SMF
	filename = pigpath.SubSpecialDirectories(filename)
	smf, err = midi.ReadSMF(filename)
	if err != nil {
		errmsg := "Can not open MIDI file %s"
		err = pigerr.CompoundError(err, fmt.Sprintf(errmsg, filename))
		return err
	}
	op.Stop()
	op.smf = smf
	return err
}

func (op *MIDIPlayer) MediaFilename() string {
	if op.smf == nil {
		return ""
	}
	return op.smf.Filename()
}

func (op *MIDIPlayer) Duration() int {
	if op.smf == nil {
		return 0.0
	}
	return int(1000 * op.smf.Duration())
}

func (op *MIDIPlayer) Position() int {
	return op.currentTime
}

func (op *MIDIPlayer) playloop() {
	time.Sleep(time.Duration(op.delayStart) * time.Millisecond)
	fmt.Printf("MIDIPlayer %s playing\n", op.Name())
	var err error
	var track *midi.SMFTrack
	track, err = op.smf.Track(0)
	if err != nil {
		return
	}
	events := track.Events()
	op.noteQueue.Reset()
	for op.isPlaying && op.eventIndex < len(events) {
		event := events[op.eventIndex]
		delay := time.Duration((op.tickDuration * float64(event.DeltaTime())) * 1000) // msec
		time.Sleep(delay * time.Millisecond)
		fmt.Printf("event [%3d]  δ %4d  delay = %4d msec  running %6d  %s\n", int(event.DeltaTime()), op.eventIndex, delay, op.currentTime, event)
		switch {
		case event.IsChannelEvent():
			pe := event.PortmidiEvent()
			op.distribute(pe)
			op.noteQueue.Update(&pe)
		case event.IsSystemEvent():
			pe := event.PortmidiEvent()
			op.distribute(pe)
		default: // meta
			mtype := byte(event.MetaType())
			switch {
			case mtype == byte(midi.META_TEMPO):
				bpm, _ := event.MetaTempoBPM()
				op.tempo = bpm
				fmt.Printf("tempo = %f\n", op.tempo)
			case midi.IsMetaTextType(mtype):
				bytes, err := event.MetaData()
				if err == nil {
					fmt.Printf("time %8d  %s : %s\n", op.currentTime, midi.MetaType(mtype), string(bytes))
				}
			case mtype == byte(midi.META_END_OF_TRACK):
				op.isPlaying = false
			default:
				// ignore
			}
		}
		op.currentTime += int(delay)
		op.eventIndex++
	}
	op.isPlaying = false
	fmt.Printf("MIDIPlayer %s play-stop\n", op.Name())
	op.killActiveNotes()
	return
}

	

	

			
		

		
	
